package store

import "github.com/bryack/lgwt_app/domain"

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() (domain.League, error)
}
