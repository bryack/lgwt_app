package database

import (
	"io"

	"github.com/bryack/lgwt_app/domain"
)

type FileSystemPlayerStore struct {
	database io.ReadSeeker
}

func (f *FileSystemPlayerStore) GetLeague() []domain.Player {
	f.database.Seek(0, io.SeekStart)
	league, _ := NewLeague(f.database)
	return league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	var wins int

	for _, v := range f.GetLeague() {
		if v.Name == name {
			wins = v.Wins
			break
		}
	}
	return wins
}
