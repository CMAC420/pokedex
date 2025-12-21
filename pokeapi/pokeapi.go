package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"pokedex/internal/pokecache"
)

var cache = pokecache.NewCache(5 * time.Second)

const baseURL = "https://pokeapi.co/api/v2"

type LocationAreaResponse struct {
	Count    int     `json:"count"`
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

	if data, ok := cache.Get(endpoint); ok {
		fmt.Println("(cache hit)")
		var areas LocationAreaResponse
		if err := json.Unmarshal(data, &areas); err != nil {
			return LocationAreaResponse{}, err
		}
		return areas, nil
	}

	fmt.Println("(cache miss)")
	resp, err := http.Get(endpoint)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return LocationAreaResponse{}, fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	cache.Add(endpoint, body)

	var areas LocationAreaResponse
	if err := json.Unmarshal(body, &areas); err != nil {
		return LocationAreaResponse{}, err
	}

	return areas, nil
}

type LocationAreaDetails struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func GetLocationAreaDetails(name string) (LocationAreaDetails, error) {
	url := baseURL + "/location-area/" + name

	if data, ok := cache.Get(url); ok {
		var details LocationAreaDetails
		if err := json.Unmarshal(data, &details); err != nil {
			return LocationAreaDetails{}, err
		}
		return details, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return LocationAreaDetails{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return LocationAreaDetails{}, fmt.Errorf("location area not found")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaDetails{}, err
	}

	cache.Add(url, body)

	var details LocationAreaDetails
	if err := json.Unmarshal(body, &details); err != nil {
		return LocationAreaDetails{}, err
	}

	return details, nil

}
