package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/golang/groupcache"
)

var ports = []string{"9091", "9092", "9093"}
var url string = "https://pokeapi.co/api/v2/pokemon/"

func getPokemon(name string) ([]byte, error) {
	response, err := http.Get(url + name)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(response.Body)
}

func makeHandler(pokemonGroup *groupcache.Group) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		query := r.URL.Query()
		values, found := query["name"]
		if !found || len(values) == 0 {
			fmt.Println("query params not found")
			os.Exit(1)
		}

		name := values[0]
		if len(name) == 0 {
			fmt.Println("no pokemon specified...")
			os.Exit(1)
		}

		var pokemonData []byte
		err := pokemonGroup.Get(r.Context(), name, groupcache.AllocatingByteSliceSink(&pokemonData))

		if err != nil {
			os.Exit(1)
		}

		w.Write(pokemonData)
	}
}

func main() {
	var instanceIdx int
	flag.IntVar(&instanceIdx, "index", 0, "groupcache index")
	flag.Parse()
	peers, pokemonGroup := InitializeCache(instanceIdx)
	mux := http.NewServeMux()

	mux.Handle("/_groupcache/", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		peers.ServeHTTP(rw, r)
		fmt.Printf("%+v\n", pokemonGroup.CacheStats(groupcache.MainCache))

	}))
	mux.Handle("/pokemon", http.HandlerFunc(makeHandler(pokemonGroup)))

	http.ListenAndServe(":"+ports[instanceIdx], mux)
}
