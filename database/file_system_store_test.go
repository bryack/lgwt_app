package database

import (
	"strings"
	"testing"

	"github.com/bryack/lgwt_app/domain"
	"github.com/stretchr/testify/assert"
)

func TestFileSystemStore(t *testing.T) {

	t.Run("league from a reader", func(t *testing.T) {
		database := strings.NewReader(`[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

		store := FileSystemPlayerStore{database: database}

		got := store.GetLeague()

		want := []domain.Player{
			{"Cleo", 10},
			{"Chris", 33},
		}

		assert.Equal(t, want, got)
	})
}
