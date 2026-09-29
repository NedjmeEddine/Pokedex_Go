package pokeapi

import (
	"net/http"
	"time"

	"github.com/NedjmeEddine/Pokedex_Go/internal/pokecach"
)

const baseURL = "https://pokeapi.co/api/v2/"

type Client struct {
	httpClient http.Client
	cache      pokecach.Cache
}

func NewClient() Client {
	return Client{
		httpClient: http.Client{
			Timeout: time.Minute,
		},
		cache: pokecach.NewCache(5 * time.Minute),
	}

}
