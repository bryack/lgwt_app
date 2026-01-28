package filesystem

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/bryack/lgwt_app/domain"
)

type FileSystemPlayerStore struct {
	database *json.Encoder
	league   domain.League
}

func NewFileSystemPlayerStore(database *os.File) (*FileSystemPlayerStore, error) {
	err := initialisePlayerDBFile(database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialise player db file %q: %w", database.Name(), err)
	}

	league, err := domain.NewLeague(database)
	if err != nil {
		return nil, fmt.Errorf("failed to load player store to file %q: %w", database.Name(), err)
	}

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{database}),
		league:   league,
	}, nil
}

func initialisePlayerDBFile(database *os.File) error {
	database.Seek(0, io.SeekStart)

	info, err := database.Stat()
	if err != nil {
		return fmt.Errorf("failed to get info of file %q: %w", database.Name(), err)
	}

	if info.Size() == 0 {
		database.Write([]byte("[]"))
		database.Seek(0, io.SeekStart)
	}

	return nil
}

func (f *FileSystemPlayerStore) GetLeague() (domain.League, error) {
	sort.Slice(f.league, func(i, j int) bool {
		return f.league[i].Wins > f.league[j].Wins
	})
	return f.league, nil
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	player := f.league.Find(name)
	if player != nil {
		return player.Wins
	}
	return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)
	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, domain.Player{Name: name, Wins: 1})
	}

	f.database.Encode(&f.league)
}
