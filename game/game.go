package game

import (
	"time"

	"github.com/bryack/lgwt_app/scheduler"
)

type Game struct {
	alerter scheduler.BlindAlerter
}

func NewGame(alerter scheduler.BlindAlerter) *Game {
	return &Game{
		alerter: alerter,
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
