package query

import (
	"github.com/amonull/rengal/filesystem"
	"github.com/amonull/rengal/where"
	"github.com/metafates/gache"
)

type queryRecord struct {
	Rank  int    `json:"rank"`
	Query string `json:"query"`
}

var cacher = gache.New[map[string]*queryRecord](
	&gache.Options{
		Path:       where.Queries(),
		FileSystem: &filesystem.GacheFs{},
	},
)
