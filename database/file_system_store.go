package database

import (
	"io"

	"github.com/bryack/lgwt_app/domain"
)

type FileSystemPlayerStore struct {
	database io.Reader
}

func (f *FileSystemPlayerStore) GetLeague() []domain.Player {
	league, _ := NewLeague(f.database)
	return league
}
