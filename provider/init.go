package provider

import (
	"github.com/amonull/rengal/provider/generic"
	"github.com/amonull/rengal/provider/mangadex"
	"github.com/amonull/rengal/provider/manganato"
	"github.com/amonull/rengal/provider/manganelo"
	"github.com/amonull/rengal/provider/mangapill"
	"github.com/amonull/rengal/source"
)

const CustomProviderExtension = ".lua"

var builtinProviders = []*Provider{
	{
		ID:   mangadex.ID,
		Name: mangadex.Name,
		CreateSource: func() (source.Source, error) {
			return mangadex.New(), nil
		},
	},
}

func init() {
	for _, conf := range []*generic.Configuration{
		manganelo.Config,
		manganato.Config,
		mangapill.Config,
	} {
		builtinProviders = append(builtinProviders, &Provider{
			ID:   conf.ID(),
			Name: conf.Name,
			CreateSource: func() (source.Source, error) {
				return generic.New(conf), nil
			},
		})
	}
}
