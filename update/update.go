package update

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/amonull/rengal/constant"
	"github.com/amonull/rengal/converter/cbz"
	"github.com/amonull/rengal/filesystem"
	"github.com/amonull/rengal/log"
	"github.com/amonull/rengal/source"
	"github.com/amonull/rengal/util"
)

func Metadata(mangaPath string) error {
	log.Infof("extracting series name from %s", mangaPath)
	name, err := GetName(mangaPath)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("extracted name: %s", name)
	log.Infof("finding %s on anilist", name)
	manga := &source.Manga{
		Name: name,
	}

	// will set new metadata from anilist
	err = manga.PopulateMetadata(func(string) {})
	if err != nil {
		log.Error()
		return err
	}

	chapters, err := getChapters(mangaPath)
	if err != nil {
		log.Error(err)
		return err
	}

	manga.Chapters = make([]*source.Chapter, 0, 0)
	chaptersPaths := make(map[*source.Chapter]string)
	prepareToUpdateComicXML(chapters, manga, chaptersPaths)

	// okay, we're ready to regenerate series.json and ComicInfo.xml now
	seriesJSON := manga.SeriesJSON()
	buf, err := json.Marshal(seriesJSON)
	if err != nil {
		log.Error(err)
		return err
	}

	// update series.json
	log.Info("updating series json")
	err = filesystem.Api().WriteFile(filepath.Join(mangaPath, "series.json"), buf, os.ModePerm)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("downloading new cover")
	updateCovers(manga, mangaPath)

	log.Infof("updating ComicInfo.xml for %d chapters", len(manga.Chapters))
	return updateComicXML(manga, chaptersPaths)
}

func prepareToUpdateComicXML(chapters []*downloadedChapter, manga *source.Manga,
	chaptersPaths map[*source.Chapter]string) {
	// since we are trying to update ComicInfo.xml here, we do not care about any other formats other than FormatCBZ
	// assuming all other formats are CBZ if first is to not run through th loop
	// also skip if no chapters in slice
	if len(chapters) < 0 || chapters[0].format != constant.FormatCBZ {
		return
	}

	for _, chapter := range chapters {
		if chapter.format != constant.FormatCBZ {
			continue
		}

		log.Infof("getting ComicInfoXML from %s", chapter.path)
		comicInfo, err := getComicInfoXML(chapter.path)
		if err != nil {
			log.Error(err)
			continue
		}

		chap := &source.Chapter{
			Name:  comicInfo.Title,
			Manga: manga,
			URL:   comicInfo.Web,
			Index: uint16(comicInfo.Number),
		}
		manga.Chapters = append(manga.Chapters, chap)
		chaptersPaths[chap] = chapter.path
	}
}

func updateCovers(manga *source.Manga, mangaPath string) {
	// remove old cover(s).
	// even though DownloadCover() will overwrite previous one
	// there may be a sitation when new cover has a different extension
	// which would result having duplicates
	files, err := filesystem.Api().ReadDir(mangaPath)
	if err == nil {
		for _, file := range files {
			if util.FileStem(file.Name()) == "cover" {
				_ = filesystem.Api().Remove(filepath.Join(mangaPath, file.Name()))
			}
		}
	}
	err = manga.DownloadCover(true, mangaPath, func(string) {})
	if err != nil {
		log.Error(err)
	}
}

func updateComicXML(manga *source.Manga, chaptersPaths map[*source.Chapter]string) error {
	for _, chapter := range manga.Chapters {
		path := chaptersPaths[chapter]
		file, err := filesystem.Api().Open(path)
		if err != nil {
			log.Error(err)
			continue
		}

		stat, err := file.Stat()
		if err != nil {
			_ = file.Close()
			continue
		}

		// go to memmap fs to unzip
		filesystem.SetMemMapFs()
		err = util.Unzip(file, stat.Size(), chapter.Name)
		if err != nil {
			log.Error(err)
			_ = file.Close()
			continue
		}

		// add pages before converting back to cbz
		files, err := filesystem.Api().ReadDir(chapter.Name)
		if err != nil {
			log.Error(err)
			_ = file.Close()
			continue
		}

		for _, file := range files {
			if err := prepareChapter(file, chapter); err != nil {
				return err
			}
		}

		_ = file.Close()

		filesystem.SetOsFs()

		if err := removeOldFile(path); err != nil {
			log.Error(err)
			continue
		}

		if err := saveChapter(path, chapter); err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}

func prepareChapter(file os.FileInfo, chapter *source.Chapter) error {
	// skip ComicInfo.xml (non-fatal no error)
	if strings.HasSuffix(file.Name(), ".xml") {
		return nil
	}

	image, err := filesystem.Api().ReadFile(filepath.Join(chapter.Name, file.Name()))
	// we can not let some pages be gone
	// so if we can't open any - whole process should stop
	if err != nil {
		log.Error(err)
		return err
	}

	chapter.Pages = append(chapter.Pages, &source.Page{
		Chapter:   chapter,
		Size:      uint64(file.Size()),
		Index:     uint16(len(chapter.Pages)),
		Extension: filepath.Ext(file.Name()),
		Contents:  bytes.NewBuffer(image),
	})

	return nil
}

func removeOldFile(path string) error {
	log.Debugf("removing old %s", path)
	return filesystem.Api().Remove(path)
}

func saveChapter(path string, chapter *source.Chapter) error {
	log.Debugf("saving to %s", path)
	return cbz.SaveTo(chapter, path)
}
