package pdf

import (
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/log"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/spf13/viper"

	"github.com/amonull/rengal/filesystem"
	"github.com/amonull/rengal/key"
	"github.com/amonull/rengal/source"
	"github.com/amonull/rengal/util"
)

type PDF struct{}

func New() *PDF {
	return &PDF{}
}

func (*PDF) Save(chapter *source.Chapter) (string, error) {
	return save(chapter, false)
}

func (*PDF) SaveTemp(chapter *source.Chapter) (string, error) {
	return save(chapter, true)
}

func save(chapter *source.Chapter, temp bool) (path string, err error) {
	path, err = chapter.Path(temp)
	if err != nil {
		return "", err
	}

	file, err := filesystem.Api().Create(path)
	if err != nil {
		return "", err
	}

	defer util.Ignore(file.Close)

	err = pagesToPDF(file, chapter.Pages)
	return "", err
}

// pagesToPDF will convert images to PDF and write to w
func pagesToPDF(w io.Writer, pages []*source.Page) error {
	conf := model.NewDefaultConfiguration()
	conf.Cmd = model.IMPORTIMAGES
	imp := pdfcpu.DefaultImportConfig()

	var (
		ctx *model.Context
		err error
	)

	ctx, err = pdfcpu.CreateContextWithXRefTable(conf, imp.PageDim)
	if err != nil {
		return err
	}

	pagesIndRef, err := ctx.Pages()
	if err != nil {
		return err
	}

	// This is the page tree root.
	pagesDict, err := ctx.DereferenceDict(*pagesIndRef)
	if err != nil {
		return err
	}

	for _, page := range pages {
		indRefs, err := pdfcpu.NewPagesForImage(ctx.XRefTable, page, pagesIndRef, imp)

		if err != nil {
			if viper.GetBool(key.FormatsSkipUnsupportedImages) {
				continue
			}

			return err
		}

		// TODO: this is a shit fix but doesnt seem to break building pdfs from images
		if err = model.AppendPageTree(indRefs[0], 1, pagesDict); err != nil {
			return err
		}

		ctx.PageCount++
	}

	if err = api.WriteContext(ctx, w); err != nil {
		return err
	}

	log.Stats.Printf("XRefTable:\n%s\n", ctx)

	return nil
}
