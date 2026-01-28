package filesystem

import (
	"testing"

	"github.com/bryack/lgwt_app/domain"
	"github.com/bryack/lgwt_app/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		database, cleanDatabase := testhelpers.CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store := NewFileSystemPlayerStore(database)
		got, err := store.GetLeague()
		assert.NoError(t, err)

		want := domain.League{
			{"Cleo", 10},
			{"Chris", 33},
		}

		assert.Equal(t, want, got)
		got, err = store.GetLeague()
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := testhelpers.CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store := NewFileSystemPlayerStore(database)

		got := store.GetPlayerScore("Chris")
		want := 33

		assert.Equal(t, want, got)
	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := testhelpers.CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store := NewFileSystemPlayerStore(database)

		store.RecordWin("Chris")
		got := store.GetPlayerScore("Chris")
		want := 34
		assert.Equal(t, want, got)
	})
	t.Run("store wins for new player", func(t *testing.T) {
		database, cleanDatabase := testhelpers.CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store := NewFileSystemPlayerStore(database)

		store.RecordWin("Pepper")
		got := store.GetPlayerScore("Pepper")
		want := 1
		assert.Equal(t, want, got)
	})
}
