package downloader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"

	"github.com/amonull/rengal/color"
	"github.com/amonull/rengal/converter"
	"github.com/amonull/rengal/filesystem"
	"github.com/amonull/rengal/history"
	"github.com/amonull/rengal/key"
	"github.com/amonull/rengal/log"
	"github.com/amonull/rengal/source"
	"github.com/amonull/rengal/style"
)

// Download the chapter using given source.
func Download(chapter *source.Chapter, progress func(string)) (string, error) {
	log.Info("downloading " + chapter.Name)

	path, err := chapter.Path(false)
	if err != nil {
		return "", err
	}

	if !viper.GetBool(key.DownloaderRedownloadExisting) && chapter.IsDownloaded() {
		log.Info("chapter already downloaded, skipping")
		return path, nil
	}

	// non-fatal err
	log.Info("chapter already downloaded, deleting and redownloading")
	err = filesystem.Api().Remove(path)
	if err != nil {
		log.Warn(err)
	}

	// fatal err
	pages, err := downloadPages(chapter, progress)
	if err != nil {
		log.Error(err)
		return "", err
	}

	// non-fatal err
	if viper.GetBool(key.MetadataFetchAnilist) {
		if err := fetchMetadata(chapter, progress); err != nil {
			log.Warn(err)
		}
	}

	// non-fatal err
	if viper.GetBool(key.MetadataSeriesJSON) {
		if err := writeMetadata(chapter, progress); err != nil {
			log.Warn(err)
		}
	}

	// non-fatal err
	if viper.GetBool(key.DownloaderDownloadCover) {
		if err := downloadCover(chapter, progress); err != nil {
			log.Warn(err)
		}
	}

	// fatal err
	if path, err = convertDownloadedContent(chapter, pages, progress); err != nil {
		log.Error(err)
		return "", err
	}

	// non-fatal err
	if viper.GetBool(key.HistorySaveOnDownload) {
		if err := saveHistory(chapter); err != nil {
			log.Warn(err)
		} else {
			log.Info("history saved")
		}
	}

	log.Info("downloaded without errors")
	progress("Downloaded")
	return path, nil
}

func downloadPages(chapter *source.Chapter, progress func(string)) ([]*source.Page, error) {
	progress("Getting pages")

	pages, err := chapter.Source().PagesOf(chapter)
	if err != nil {
		return nil, err
	}
	log.Info("found " + strconv.Itoa(len(pages)) + " pages")

	err = chapter.DownloadPages(false, progress)
	if err != nil {
		return nil, err
	}

	return pages, nil
}

// uses anilist to fetch metadata
func fetchMetadata(chapter *source.Chapter, progress func(string)) error {
	return chapter.Manga.PopulateMetadata(progress)
}

// writes metadata as json
func writeMetadata(chapter *source.Chapter, progress func(string)) error {
	path, err := chapter.Manga.Path(false)
	if err != nil {
		return err
	}

	path = filepath.Join(path, "series.json")
	progress("Generating series.json")
	seriesJSON := chapter.Manga.SeriesJSON()

	buf, err := json.Marshal(seriesJSON)
	if err != nil {
		return err
	}

	return filesystem.Api().WriteFile(path, buf, os.ModePerm)
}

// downloads the cover of content
func downloadCover(chapter *source.Chapter, progress func(string)) error {
	coverDir, err := chapter.Manga.Path(false)
	if err != nil {
		return err
	}

	return chapter.Manga.DownloadCover(false, coverDir, progress)
}

// converts the downloaded content onto user specified format
func convertDownloadedContent(chapter *source.Chapter, pages []*source.Page, progress func(string)) (string, error) {
	log.Info("getting " + viper.GetString(key.FormatsUse) + " converter")
	progress(fmt.Sprintf(
		"Converting %d pages to %s %s",
		len(pages),
		style.Fg(color.Yellow)(viper.GetString(key.FormatsUse)),
		style.Faint(chapter.SizeHuman())),
	)
	conv, err := converter.Get(viper.GetString(key.FormatsUse))
	if err != nil {
		log.Error(err)
		return "", err
	}

	log.Info("converting " + viper.GetString(key.FormatsUse))
	return conv.Save(chapter)
}

func saveHistory(chapter *source.Chapter) error {
	return history.Save(chapter)
}
