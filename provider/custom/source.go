package custom

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"

	"github.com/amonull/rengal/source"
)

type luaSource struct {
	name  string
	state *lua.LState
	cache struct {
		mangas   *cacher[[]*source.Manga]
		chapters *cacher[[]*source.Chapter]
	}
}

func (s *luaSource) Name() string {
	return s.name
}

func (s *luaSource) ID() string {
	return IDfromName(s.name)
}

func newLuaSource(name string, state *lua.LState) *luaSource {
	s := &luaSource{
		name:  name,
		state: state,
	}

	cacheName := func(cacheFor string) string {
		return fmt.Sprintf("%s_%s", s.ID(), cacheFor)
	}

	s.cache.mangas = newCacher[[]*source.Manga](cacheName("mangas"))
	s.cache.chapters = newCacher[[]*source.Chapter](cacheName("chapters"))

	return s
}

//nolint:unparam // lua.LValue ret is infact not being used ever but it will be kept to not break functionality that might depend on black magic
func (s *luaSource) call(fn string, ret lua.LValueType, args ...lua.LValue) (lua.LValue, error) {
	err := s.state.CallByParam(lua.P{
		Fn:      s.state.GetGlobal(fn),
		NRet:    1,
		Protect: true,
	}, args...)

	if err != nil {
		return nil, err
	}

	val := s.state.Get(-1)

	if val.Type() != ret {
		s.state.RaiseError(fn + " was expected to return a " + ret.String() + ", got " + val.Type().String())
	}

	return val, nil
}
