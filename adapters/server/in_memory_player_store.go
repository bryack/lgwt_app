package server

import (
	"sync"

	"github.com/bryack/lgwt_app/domain"
)

type InMemoryPlayerStore struct {
	store map[string]int
	lock  sync.RWMutex
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		store: map[string]int{},
		lock:  sync.RWMutex{},
	}
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.store[name]
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.store[name]++
}

func (i *InMemoryPlayerStore) GetLeague() ([]domain.Player, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()
	league := make([]domain.Player, 0, len(i.store))
	for name, wins := range i.store {
		league = append(league, domain.Player{Name: name, Wins: wins})
	}
	return league, nil
}
