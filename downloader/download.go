package downloader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/amonull/rengal/color"
	"github.com/amonull/rengal/converter"
	"github.com/amonull/rengal/filesystem"
	"github.com/amonull/rengal/history"
	"github.com/amonull/rengal/key"
	"github.com/amonull/rengal/log"
	"github.com/amonull/rengal/source"
	"github.com/amonull/rengal/style"
	"github.com/spf13/viper"
)

// Download the chapter using given source.
func Download(chapter *source.Chapter, progress func(string)) (string, error) {
	log.Info("downloading " + chapter.Name)

	path, err := chapter.Path(false)
	if err != nil {
		return "", err
	}

	if viper.GetBool(key.DownloaderRedownloadExisting) {
		log.Info("chapter already downloaded, deleting and redownloading")
		err = filesystem.Api().Remove(path)
		if err != nil {
			log.Warn(err)
		}
	} else {
		log.Info("checking if chapter is already downloaded")
		if chapter.IsDownloaded() {
			log.Info("chapter already downloaded, skipping")
			return path, nil
		}
	}

	progress("Getting pages")
	pages, err := chapter.Source().PagesOf(chapter)
	if err != nil {
		log.Error(err)
		return "", err
	}
	log.Info("found " + fmt.Sprintf("%d", len(pages)) + " pages")

	err = chapter.DownloadPages(false, progress)
	if err != nil {
		log.Error(err)
		return "", err
	}

	if viper.GetBool(key.MetadataFetchAnilist) {
		err := chapter.Manga.PopulateMetadata(progress)
		if err != nil {
			log.Warn(err)
		}
	}

	if viper.GetBool(key.MetadataSeriesJSON) {
		path, err := chapter.Manga.Path(false)
		if err != nil {
			log.Warn(err)
		} else {
			path = filepath.Join(path, "series.json")
			progress("Generating series.json")
			seriesJSON := chapter.Manga.SeriesJSON()
			buf, err := json.Marshal(seriesJSON)
			if err != nil {
				log.Warn(err)
			} else {
				err = filesystem.Api().WriteFile(path, buf, os.ModePerm)
				if err != nil {
					log.Warn(err)
				}
			}
		}
	}

	if viper.GetBool(key.DownloaderDownloadCover) {
		coverDir, err := chapter.Manga.Path(false)
		if err == nil {
			_ = chapter.Manga.DownloadCover(false, coverDir, progress)
		}
	}

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
	path, err = conv.Save(chapter)
	if err != nil {
		log.Error(err)
		return "", err
	}

	if viper.GetBool(key.HistorySaveOnDownload) {
		go func() {
			err = history.Save(chapter)
			if err != nil {
				log.Warn(err)
			} else {
				log.Info("history saved")
			}
		}()
	}

	log.Info("downloaded without errors")
	progress("Downloaded")
	return path, nil
}
