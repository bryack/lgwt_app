package database

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/bryack/lgwt_app/domain"
)

func NewLeague(rdr io.Reader) ([]domain.Player, error) {
	var league []domain.Player
	err := json.NewDecoder(rdr).Decode(&league)
	if err != nil {
		err = fmt.Errorf("failed to parse league: %w", err)
	}
	return league, nil
}
