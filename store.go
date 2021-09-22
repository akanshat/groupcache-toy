package main

import (
	"context"

	"github.com/golang/groupcache"
)

type Store interface {
	// set()
	Get(name string) ([]byte, error)
}

type DirectAPI struct {
}

var _ Store = DirectAPI{}

func (d DirectAPI) Get(name string) ([]byte, error) {
	return getPokemon(name)
}

type InMemCache struct {
	mp map[string][]byte
}

var _ Store = &InMemCache{}

func (m *InMemCache) Get(name string) ([]byte, error) {
	if _, ok := m.mp[name]; ok {
		return m.mp[name], nil
	}
	res, err := getPokemon(name)
	m.mp[name] = res
	return res, err
}

type GroupCache struct {
	group *groupcache.Group
}

var _ Store = &GroupCache{}

func (g *GroupCache) Get(name string) ([]byte, error) {
	var pokemonData []byte
	err := g.group.Get(context.Background(), name, groupcache.AllocatingByteSliceSink(&pokemonData))
	return pokemonData, err
}
