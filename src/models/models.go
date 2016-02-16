package models

import (
	"encoding/json"
)

// Movie is OMDB-sourced movie info
type Movie struct {
	Title    string
	Actors   string
	Poster   string
	Year     string
	Plot     string
	Director string
	Rated    string
	Runtime  string
	Genre    string
	Rating   string `json:"imdbRating"`
	ImdbID   string `json:"imdbID"`
	Seen     bool   `json:"seen"`
}

func (m *Movie) String() string {
	return m.Title
}

// MarshalBinary allows writes to Redis
func (m *Movie) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary allows reads from Redis
func (m *Movie) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
