package cli_test

import (
	"strings"
	"testing"

	"github.com/bryack/lgwt_app/adapters/cli"
	"github.com/bryack/lgwt_app/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestCLI(t *testing.T) {
	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &testhelpers.StubPlayerStore{}
		c := cli.NewCLI(playerStore, in)
		c.PlayPoker()
		assert.True(t, len(playerStore.WinCalls) != 0)

		got := playerStore.WinCalls[0]
		want := "Chris"
		assert.Equal(t, want, got)
	})
	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &testhelpers.StubPlayerStore{}
		c := cli.NewCLI(playerStore, in)
		c.PlayPoker()
		assert.True(t, len(playerStore.WinCalls) != 0)

		got := playerStore.WinCalls[0]
		want := "Cleo"
		assert.Equal(t, want, got)
	})
}
