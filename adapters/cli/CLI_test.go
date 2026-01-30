package cli_test

import (
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
		in := strings.NewReader("Chris wins\n")
		playerStore := &testhelpers.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}

		c := cli.NewCLI(playerStore, in, blindAlerter)
		c.PlayPoker()

		if len(blindAlerter.alerts) != 1 {
			t.Fatal("expected a blind alert to be scheduled")
		}
	})
}
