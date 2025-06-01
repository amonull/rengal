package generic

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"

	"github.com/amonull/rengal/constant"
	"github.com/amonull/rengal/source"
	"github.com/amonull/rengal/where"
)

// New generates a new scraper with given configuration
func New(conf *Configuration) source.Source {
	scraper := Scraper{
		mangas:   make(map[string][]*source.Manga),
		chapters: make(map[string][]*source.Chapter),
		pages:    make(map[string][]*source.Page),
		config:   conf,
	}

	collectorOptions := newCollectorOpts()

	baseCollector := colly.NewCollector(collectorOptions...)
	baseCollector.SetRequestTimeout(20 * time.Second)

	mangasCollector := newMangaCollector(baseCollector, scraper)

	// Get mangas
	mangaCollector_getManga(mangasCollector, &scraper)

	_ = mangasCollector.Limit(&colly.LimitRule{
		Parallelism: int(scraper.config.Parallelism),
		RandomDelay: scraper.config.Delay,
		DomainGlob:  "*",
	})

	chaptersCollector := newChapterCollector(baseCollector, scraper)

	// Get chapters
	chaptersCollector_getChapters(chaptersCollector, &scraper)

	_ = chaptersCollector.Limit(&colly.LimitRule{
		Parallelism: int(scraper.config.Parallelism),
		RandomDelay: scraper.config.Delay,
		DomainGlob:  "*",
	})

	pagesCollector := newPagesCollector(baseCollector)

	// Get pages
	pagesCollector_getPages(pagesCollector, &scraper)

	_ = pagesCollector.Limit(&colly.LimitRule{
		Parallelism: int(scraper.config.Parallelism),
		RandomDelay: scraper.config.Delay,
		DomainGlob:  "*",
	})

	scraper.mangasCollector = mangasCollector
	scraper.chaptersCollector = chaptersCollector
	scraper.pagesCollector = pagesCollector

	return &scraper
}

func newCollectorOpts() []colly.CollectorOption {
	return []colly.CollectorOption{
		colly.AllowURLRevisit(),
		colly.Async(true),
		colly.CacheDir(where.Cache()),
	}
}

func newMangaCollector(baseCollector *colly.Collector, scraper Scraper) *colly.Collector {
	collector := baseCollector.Clone()
	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", "https://google.com")
		r.Headers.Set("accept-language", "en-US")
		r.Headers.Set("Accept", "text/html")
		r.Headers.Set("Host", scraper.config.BaseURL)
		r.Headers.Set("User-Agent", constant.UserAgent)
	})

	return collector
}

func newChapterCollector(baseCollector *colly.Collector, scraper Scraper) *colly.Collector {
	collector := baseCollector.Clone()
	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", r.Ctx.GetAny("manga").(*source.Manga).URL)
		r.Headers.Set("accept-language", "en-US")
		r.Headers.Set("Accept", "text/html")
		r.Headers.Set("Host", scraper.config.BaseURL)
		r.Headers.Set("User-Agent", constant.UserAgent)
	})

	return collector
}

func newPagesCollector(baseCollector *colly.Collector) *colly.Collector {
	collector := baseCollector.Clone()
	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", r.Ctx.GetAny("chapter").(*source.Chapter).URL)
		r.Headers.Set("accept-language", "en-US")
		r.Headers.Set("Accept", "text/html")
		r.Headers.Set("User-Agent", constant.UserAgent)
	})

	return collector
}

func mangaCollector_getManga(mangasCollector *colly.Collector, scraper *Scraper) {
	mangasCollector.OnHTML("html", func(e *colly.HTMLElement) {
		elements := e.DOM.Find(scraper.config.MangaExtractor.Selector)
		path := e.Request.URL.String()
		scraper.mangas[path] = make([]*source.Manga, elements.Length())

		elements.Each(func(i int, selection *goquery.Selection) {
			link := scraper.config.MangaExtractor.URL(selection)
			url := e.Request.AbsoluteURL(link)
			manga := source.Manga{
				Name:     scraper.config.MangaExtractor.Name(selection),
				URL:      url,
				Index:    uint16(e.Index),
				Chapters: make([]*source.Chapter, 0),
				ID:       filepath.Base(url),
				Source:   scraper,
			}
			manga.Metadata.Cover.ExtraLarge = scraper.config.MangaExtractor.Cover(selection)

			scraper.mangas[path][i] = &manga
		})
	})
}

func chaptersCollector_getChapters(chaptersCollector *colly.Collector, scraper *Scraper) {
	chaptersCollector.OnHTML("html", func(e *colly.HTMLElement) {
		elements := e.DOM.Find(scraper.config.ChapterExtractor.Selector)
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		scraper.chapters[path] = make([]*source.Chapter, elements.Length())
		manga := e.Request.Ctx.GetAny("manga").(*source.Manga)

		elements.Each(func(i int, selection *goquery.Selection) {
			link := scraper.config.ChapterExtractor.URL(selection)
			url := e.Request.AbsoluteURL(link)

			chapter := source.Chapter{
				Name:   scraper.config.ChapterExtractor.Name(selection),
				URL:    url,
				Index:  uint16(e.Index),
				Pages:  make([]*source.Page, 0),
				ID:     filepath.Base(url),
				Manga:  manga,
				Volume: scraper.config.ChapterExtractor.Volume(selection),
			}
			scraper.chapters[path][i] = &chapter
		})
		manga.Chapters = scraper.chapters[path]
	})
}

func pagesCollector_getPages(pagesCollector *colly.Collector, scraper *Scraper) {
	pagesCollector.OnHTML("html", func(e *colly.HTMLElement) {
		elements := e.DOM.Find(scraper.config.PageExtractor.Selector)
		path := e.Request.AbsoluteURL(e.Request.URL.Path)
		scraper.pages[path] = make([]*source.Page, elements.Length())
		chapter := e.Request.Ctx.GetAny("chapter").(*source.Chapter)

		elements.Each(func(i int, selection *goquery.Selection) {
			link := scraper.config.PageExtractor.URL(selection)
			ext := filepath.Ext(link)
			// remove some query params from the extension
			ext = strings.Split(ext, "?")[0]

			page := source.Page{
				URL:       link,
				Index:     uint16(i),
				Chapter:   chapter,
				Extension: ext,
			}
			scraper.pages[path][i] = &page
		})
		chapter.Pages = scraper.pages[path]
	})
}
