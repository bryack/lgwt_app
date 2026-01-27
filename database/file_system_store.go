package database

import (
	"encoding/json"
	"io"

	"github.com/bryack/lgwt_app/server"
)

type FileSystemPlayerStore struct {
	database io.Reader
}

func (f *FileSystemPlayerStore) GetLeague() []server.Player {
	var league []server.Player
	json.NewDecoder(f.database).Decode(&league)
	return league
}
