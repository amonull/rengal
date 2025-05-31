package custom

import (
	"path/filepath"
	"time"

	"github.com/metafates/gache"
	"github.com/samber/mo"

	"github.com/amonull/rengal/filesystem"
	"github.com/amonull/rengal/util"
	"github.com/amonull/rengal/where"
)

type cacher[T any] struct {
	internal *gache.Cache[map[string]T]
}

func newCacher[T any](name string) *cacher[T] {
	return &cacher[T]{
		internal: gache.New[map[string]T](&gache.Options{
			Lifetime:   time.Hour * 24,
			Path:       filepath.Join(where.Cache(), util.SanitizeFilename(name)+".json"),
			FileSystem: &filesystem.GacheFs{},
		}),
	}
}

func (c *cacher[T]) Get(key string) mo.Option[T] {
	data, expired, err := c.internal.Get()
	if err != nil || expired || data == nil {
		return mo.None[T]()
	}

	if x, ok := data[key]; ok {
		return mo.Some(x)
	}

	return mo.None[T]()
}

func (c *cacher[T]) Set(key string, t T) error {
	data, expired, err := c.internal.Get()

	if err != nil {
		return err
	}

	if expired || data == nil {
		data = make(map[string]T)
	}

	data[key] = t

	return c.internal.Set(data)
}
