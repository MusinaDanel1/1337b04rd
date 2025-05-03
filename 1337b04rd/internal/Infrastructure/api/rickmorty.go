package api

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
)

// TODO: reduce the number of requests, save locally used avatars so that they are not repeated

type Character struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Response struct {
	Results []Character `json:"results"`
}

type RickMortyClient struct{}

func (c *RickMortyClient) GetRandomCharacter(exclude map[string]bool) (Character, error) {
	resp, err := http.Get("https://rickandmortyapi.com/api/character")
	if err != nil {
		return Character{}, err
	}
	defer resp.Body.Close()

	var data Response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return Character{}, err
	}

	var available []Character
	for _, char := range data.Results {
		if !exclude[char.Name] {
			available = append(available, char)
		}
	}

	if len(available) == 0 {
		return Character{}, errors.New("no unique characters left")
	}

	selected := available[rand.Intn(len(available))]
	return selected, nil
}
