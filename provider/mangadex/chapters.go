package mangadex

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/darylhjd/mangodex"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"

	"github.com/amonull/rengal/key"
	"github.com/amonull/rengal/source"
	"github.com/amonull/rengal/util"
)

func (m *Mangadex) ChaptersOf(manga *source.Manga) ([]*source.Chapter, error) {
	if cached, ok := m.cache.chapters.Get(manga.URL).Get(); ok {
		for _, chapter := range cached {
			chapter.Manga = manga
		}

		return cached, nil
	}

	params := url.Values{}
	params.Set("limit", strconv.Itoa(500))
	ratings := []string{mangodex.Safe, mangodex.Suggestive}
	for _, rating := range ratings {
		params.Add("contentRating[]", rating)
	}

	if viper.GetBool(key.MangadexNSFW) {
		params.Add("contentRating[]", mangodex.Porn)
		params.Add("contentRating[]", mangodex.Erotica)
	}

	// scanlation group for the chapter
	params.Add("includes[]", mangodex.ScanlationGroupRel)
	params.Set("order[chapter]", "asc")

	var (
		chapters   []*source.Chapter
		currOffset = 0
	)

	for {
		params.Set("offset", strconv.Itoa(currOffset))
		list, err := m.client.Chapter.GetMangaChapters(manga.ID, params)
		if err != nil {
			return nil, err
		}

		for i, chapter := range list.Data {
			preparedChapter, offset := prepareChapter(chapter, i, manga)
			currOffset += offset
			if preparedChapter == nil {
				continue
			}

			manga.Chapters = append(manga.Chapters, preparedChapter)
		}
		currOffset += 500
		if currOffset >= list.Total {
			break
		}

		if currOffset >= list.Total {
			break
		}
	}

	slices.SortFunc(chapters, func(a, b *source.Chapter) int {
		return util.CompareInt(int(a.Index), int(b.Index))
	})

	manga.Chapters = chapters
	_ = m.cache.chapters.Set(manga.URL, chapters)
	return chapters, nil
}

func prepareChapter(chapter mangodex.Chapter, index int, manga *source.Manga) (*source.Chapter, int) {
	language := viper.GetString(key.MangadexLanguage)

	// Skip external chapters. Their pages cannot be downloaded.
	if chapter.Attributes.ExternalURL != nil && !viper.GetBool(key.MangadexShowUnavailableChapters) {
		return nil, 0
	}

	// skip chapters that are not in the current language
	if language != "any" && chapter.Attributes.TranslatedLanguage != language {
		// 500 to be added to currOffset
		return nil, 500
	}

	name := chapter.GetTitle()
	if name == "" {
		name = fmt.Sprintf("Chapter %s", chapter.GetChapterNum())
	} else {
		name = fmt.Sprintf("Chapter %s - %s", chapter.GetChapterNum(), name)
	}

	var volume string
	if chapter.Attributes.Volume != nil {
		volume = fmt.Sprintf("Vol.%s", *chapter.Attributes.Volume)
	}
	return &source.Chapter{
		Name:   name,
		Index:  uint16(index),
		ID:     chapter.ID,
		URL:    fmt.Sprintf("https://mangadex.org/chapter/%s", chapter.ID),
		Manga:  manga,
		Volume: volume,
	}, 0
}
