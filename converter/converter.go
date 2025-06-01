package converter

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"github.com/amonull/rengal/constant"
	"github.com/amonull/rengal/converter/cbz"
	"github.com/amonull/rengal/converter/pdf"
	"github.com/amonull/rengal/converter/plain"
	"github.com/amonull/rengal/converter/zip"
	"github.com/amonull/rengal/source"
)

// Converter is the interface that all converters must implement.
type Converter interface {
	Save(chapter *source.Chapter) (string, error)
	SaveTemp(chapter *source.Chapter) (string, error)
}

var converters = map[string]Converter{
	constant.FormatPlain: plain.New(),
	constant.FormatCBZ:   cbz.New(),
	constant.FormatPDF:   pdf.New(),
	constant.FormatZIP:   zip.New(),
}

// Available returns a list of available converters.
func Available() []string {
	return lo.Keys(converters)
}

// Get returns a converter by name.
// If the converter is not available, an error is returned.
func Get(name string) (Converter, error) {
	if converter, ok := converters[name]; ok {
		return converter, nil
	}

	return nil, fmt.Errorf("unkown format \"%s\", available options are %s", name, strings.Join(Available(), ", "))
}
