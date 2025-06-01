package mangadex

import (
	"path/filepath"
	"time"

	"github.com/metafates/gache"
	"github.com/samber/mo"

	"github.com/amonull/rengal/filesystem"
	"github.com/amonull/rengal/where"
)

type cacher[T any] struct {
	internal *gache.Cache[map[string]T]
}

func newCacher[T any](name string) *cacher[T] {
	return &cacher[T]{
		internal: gache.New[map[string]T](
			&gache.Options{
				Path:       filepath.Join(where.Cache(), name+".json"),
				Lifetime:   time.Hour * 24,
				FileSystem: &filesystem.GacheFs{},
			},
		),
	}
}

func (c *cacher[T]) Get(key string) mo.Option[T] {
	cached, expired, err := c.internal.Get()
	if err != nil || expired || cached == nil {
		return mo.None[T]()
	}

	if value, ok := cached[key]; ok {
		return mo.Some[T](value)
	}

	return mo.None[T]()
}

func (c *cacher[T]) Set(key string, value T) error {
	cached, expired, err := c.internal.Get()
	if err != nil {
		return err
	}

	if expired || cached == nil {
		cached = map[string]T{}
	}

	cached[key] = value
	return c.internal.Set(cached)
}
