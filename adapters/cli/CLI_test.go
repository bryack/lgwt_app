package cli_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/bryack/lgwt_app/adapters/cli"
	"github.com/bryack/lgwt_app/testhelpers"
	"github.com/stretchr/testify/assert"
)

var dummySpyAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &testhelpers.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

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

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{at: at, amount: amount})
}

func TestCLI(t *testing.T) {
	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &testhelpers.StubPlayerStore{}
		c := cli.NewCLI(playerStore, in, dummyStdOut, dummySpyAlerter)
		c.PlayPoker()
		assert.True(t, len(playerStore.WinCalls) != 0)

		got := playerStore.WinCalls[0]
		want := "Chris"
		assert.Equal(t, want, got)
	})
	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &testhelpers.StubPlayerStore{}
		c := cli.NewCLI(playerStore, in, dummyStdOut, dummySpyAlerter)
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
		c := cli.NewCLI(playerStore, in, dummyStdOut, blindAlerter)
		c.PlayPoker()

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
				if len(blindAlerter.alerts) <= 1 {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				alert := blindAlerter.alerts[i]

				assert.Equal(t, alert.amount, tt.amount)
				assert.Equal(t, alert.at, tt.at)
			})
		}
	})

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		c := cli.NewCLI(dummyPlayerStore, dummyStdIn, stdout, dummySpyAlerter)
		c.PlayPoker()

		got := stdout.String()
		want := cli.PlayerPrompt

		assert.Equal(t, want, got)
	})
}
