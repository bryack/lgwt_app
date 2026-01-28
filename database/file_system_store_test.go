package database

import (
	"strings"
	"testing"

	"github.com/bryack/lgwt_app/domain"
	"github.com/stretchr/testify/assert"
)

func TestFileSystemStore(t *testing.T) {
	database := strings.NewReader(`[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

	store := FileSystemPlayerStore{database: database}

	t.Run("league from a reader", func(t *testing.T) {
		got := store.GetLeague()

		want := []domain.Player{
			{"Cleo", 10},
			{"Chris", 33},
		}

		assert.Equal(t, want, got)
		got = store.GetLeague()
		assert.Equal(t, want, got)
	})

	t.Run("get player score", func(t *testing.T) {
		got := store.GetPlayerScore("Chris")
		want := 33

		assert.Equal(t, want, got)
	})
}
