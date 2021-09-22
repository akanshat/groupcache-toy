package main

import (
	"context"

	"github.com/golang/groupcache"
)

var pool = []string{"http://127.0.0.1:9091", "http://127.0.0.1:9092", "http://127.0.0.1:9093"}

func InitializeCache(instanceIdx int) (*groupcache.HTTPPool, *groupcache.Group) {

	currInstance := pool[instanceIdx]
	peers := groupcache.NewHTTPPool(currInstance)
	peers.Set(pool...)

	return peers, groupcache.NewGroup("pokemons", 64<<20, groupcache.GetterFunc(func(ctx context.Context, key string, dest groupcache.Sink) error {
		pokemon, err := getPokemon(key)
		if err != nil {
			return err
		}
		dest.SetBytes(pokemon)
		return nil
	}))
}
