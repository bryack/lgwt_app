package game_test

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/bryack/lgwt_app/game"
	"github.com/bryack/lgwt_app/testhelpers"
	"github.com/stretchr/testify/assert"
)

type scheduledAlert struct {
	at     time.Duration
	amount int
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int, to io.Writer) {
	s.alerts = append(s.alerts, scheduledAlert{at: at, amount: amount})
}

func TestGame_Start(t *testing.T) {
	t.Run("it schedules printing of blind values", func(t *testing.T) {
		store := &testhelpers.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}
		g := game.NewGame(blindAlerter, store)

		g.Start(5, io.Discard)

		tests := []scheduledAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, tt := range tests {
			t.Run(fmt.Sprint(tt), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				alert := blindAlerter.alerts[i]

				assert.Equal(t, alert.amount, tt.amount)
				assert.Equal(t, alert.at, tt.at)
			})
		}
	})
}

func TestGame_Finish(t *testing.T) {
	store := &testhelpers.StubPlayerStore{}
	blindAlerter := &SpyBlindAlerter{}
	g := game.NewGame(blindAlerter, store)

	winner := "Cleo"
	g.Finish(winner)

	assert.Equal(t, 1, len(store.WinCalls))
	assert.Equal(t, winner, store.WinCalls[0])
}
