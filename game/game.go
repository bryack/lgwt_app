package game

import (
	"time"

	"github.com/bryack/lgwt_app/scheduler"
	"github.com/bryack/lgwt_app/store"
)

type Game struct {
	alerter scheduler.BlindAlerter
	store   store.PlayerStore
}

func NewGame(alerter scheduler.BlindAlerter, store store.PlayerStore) *Game {
	return &Game{
		alerter: alerter,
		store:   store,
	}
}

func (g *Game) ScheduleBlindAlerts(numberOfPlayers int) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		g.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + blindIncrement
	}
}

func (g *Game) Start(numberOfPlayers int) {
	g.ScheduleBlindAlerts(numberOfPlayers)
}

func (g *Game) Finish(winner string) {
	g.store.RecordWin(winner)
}
