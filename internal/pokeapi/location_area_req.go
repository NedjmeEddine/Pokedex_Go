package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetLocationAreas(url *string) (*LocationAreasRes, error) {
	requestURL := baseURL + "location-area/"
	if url != nil {
		requestURL = *url
	}
	//cache
	data, ok := c.cache.Get(requestURL)
	if ok {
		fmt.Println("Cache hit for URL:", requestURL)
		var locationareasRes LocationAreasRes
		err := json.Unmarshal(data, &locationareasRes)
		if err != nil {
			return nil, err
		}
		return &locationareasRes, nil
	}
	fmt.Println("Cache miss for URL:", requestURL)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var locationareasRes LocationAreasRes
	err = json.Unmarshal(dat, &locationareasRes)
	if err != nil {
		return nil, err
	}
	c.cache.Add(requestURL, dat)
	return &locationareasRes, nil
}

func (c *Client) GetArea(area *string) (*Area, error) {
	requestURL := baseURL + "location-area/" + *area
	//cache
	data, ok := c.cache.Get(requestURL)
	if ok {
		fmt.Println("Cache hit for URL:", requestURL)
		var Area Area
		err := json.Unmarshal(data, &Area)
		if err != nil {
			return nil, err
		}
		return &Area, nil
	}
	fmt.Println("Cache miss for URL:", requestURL)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var Area Area
	err = json.Unmarshal(dat, &Area)
	if err != nil {
		return nil, err
	}
	c.cache.Add(requestURL, dat)
	return &Area, nil
}

func (c *Client) GetPokemonArea(pokemon *string) (*PokemonLocation, error) {
	requestURL := baseURL + "pokemon/" + *pokemon + "/encounters"
	//cache
	data, ok := c.cache.Get(requestURL)
	if ok {
		fmt.Println("Cache hit for URL:", requestURL)
		var PokemonLocation PokemonLocation
		err := json.Unmarshal(data, &PokemonLocation)
		if err != nil {
			return nil, err
		}
		return &PokemonLocation, nil
	}
	fmt.Println("Cache miss for URL:", requestURL)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var PokemonLocation PokemonLocation
	err = json.Unmarshal(dat, &PokemonLocation)
	if err != nil {
		return nil, err
	}
	c.cache.Add(requestURL, dat)
	return &PokemonLocation, nil
}
