package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "https://pokeapi.co/api/v2"

type LocationAreaResponse struct {
	Count    int     `json:"count:"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocationAreas(url *string) (LocationAreaResponse, error) {
	endpoint := baseURL + "/location-area"

	if url != nil {
		endpoint = *url
	}

	resp, err := http.Get(endpoint)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return LocationAreaResponse{}, fmt.Errorf("bad status: %s", resp.Status)
	}

	var data LocationAreaResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return LocationAreaResponse{}, err
	}

	return data, nil
}
