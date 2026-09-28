package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetLocationAreas(url *string) (*LocationAreaRes, error) {
	requestURL := baseURL + "location-area/"
	if url != nil {
		requestURL = *url
	}
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
	var locationAreaRes LocationAreaRes
	err = json.Unmarshal(dat, &locationAreaRes)
	if err != nil {
		return nil, err
	}
	return &locationAreaRes, nil
}
