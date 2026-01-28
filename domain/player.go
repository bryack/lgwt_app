package domain

import (
	"encoding/json"
	"fmt"
	"io"
)

type Player struct {
	Name string `json:"name"`
	Wins int    `json:"wins"`
}

type League []Player

func NewLeague(rdr io.Reader) ([]Player, error) {
	var league []Player
	err := json.NewDecoder(rdr).Decode(&league)
	if err != nil {
		return nil, fmt.Errorf("failed to parse league: %w", err)
	}
	return league, nil
}

func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}
