package cli_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/bryack/lgwt_app/adapters/cli"
	"github.com/bryack/lgwt_app/testhelpers"
	"github.com/stretchr/testify/assert"
)

var dummySpyAlerter = &SpyBlindAlerter{}

type SpyBlindAlerter struct {
	alerts []struct {
		scheduledAt time.Duration
		amount      int
	}
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, struct {
		scheduledAt time.Duration
		amount      int
	}{duration, amount})
}

func TestCLI(t *testing.T) {
	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &testhelpers.StubPlayerStore{}
		c := cli.NewCLI(playerStore, in, dummySpyAlerter)
		c.PlayPoker()
		assert.True(t, len(playerStore.WinCalls) != 0)

		got := playerStore.WinCalls[0]
		want := "Chris"
		assert.Equal(t, want, got)
	})
	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &testhelpers.StubPlayerStore{}
		c := cli.NewCLI(playerStore, in, dummySpyAlerter)
		c.PlayPoker()
		assert.True(t, len(playerStore.WinCalls) != 0)

		got := playerStore.WinCalls[0]
		want := "Cleo"
		assert.Equal(t, want, got)
	})

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &testhelpers.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}
		c := cli.NewCLI(playerStore, in, blindAlerter)
		c.PlayPoker()

		tests := []struct {
			expectedScheduleTime time.Duration
			expectedAmount       int
		}{
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
			t.Run(fmt.Sprintf("%d scheduled for %v", tt.expectedAmount, tt.expectedScheduleTime), func(t *testing.T) {
				if len(blindAlerter.alerts) <= 1 {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				alert := blindAlerter.alerts[i]

				amountGot := alert.amount
				if amountGot != tt.expectedAmount {
					t.Errorf("got amount %d, want %d", amountGot, tt.expectedAmount)
				}

				gotScheduledTime := alert.scheduledAt
				if gotScheduledTime != tt.expectedScheduleTime {
					t.Errorf("got scheduled time of %v, want %v", gotScheduledTime, tt.expectedScheduleTime)
				}
			})
		}
	})
}
