package main

import (
	"encoding/json"
	"net/http"
)

type countAPIConfig struct {
	Key string
}

type countAPIResult struct {
	Value int `json:"value"`
}

func (c countAPIConfig) Hit() (int, error) {
	resp, err := http.Get("https://api.countapi.xyz/hit/" + c.Key)
	if err != nil {
		return 0, err
	}

	var result countAPIResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}

	return result.Value, nil
}
