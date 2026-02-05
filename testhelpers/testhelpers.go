package testhelpers

import (
	"os"
	"testing"

	"github.com/bryack/lgwt_app/domain"
)

type StubPlayerStore struct {
	Scores   map[string]int
	WinCalls []string
	League   []domain.Player
	Err      error
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.Scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.WinCalls = append(s.WinCalls, name)
}

func (s *StubPlayerStore) GetLeague() (domain.League, error) {
	if s.Err != nil {
		return domain.League{}, s.Err
	}
	return s.League, nil
}

func AssertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.WinCalls) <= 0 {
		t.Fatalf("length of calls to RecordWin should be at least 1, got %d", len(store.WinCalls))
	}

	if store.WinCalls[0] != winner {
		t.Errorf("expected winner is %q, got %q", winner, store.WinCalls[0])
	}
}

func CreateTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	tmpFile.Write([]byte(initialData))

	removeFile := func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	return tmpFile, removeFile
}
