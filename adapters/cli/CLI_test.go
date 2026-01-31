package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bryack/lgwt_app/adapters/cli"
	"github.com/bryack/lgwt_app/testhelpers"
	"github.com/stretchr/testify/assert"
)

type SpyGame struct {
	startCalledWith int
}

func (s *SpyGame) Start(numberOfPlayers int) {
	s.startCalledWith = numberOfPlayers
}

func TestCLI(t *testing.T) {
	t.Run("it prompts the user to enter the number of players and starts the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		game := &SpyGame{}

		playerStore := &testhelpers.StubPlayerStore{}

		c := cli.NewCLI(playerStore, in, stdout, game)
		c.PlayPoker()

		assert.Equal(t, cli.PlayerPrompt, stdout.String())
		assert.Equal(t, 7, game.startCalledWith)
	})
	t.Run("record chris win from user input", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("1\nChris wins\n")
		playerStore := &testhelpers.StubPlayerStore{}
		game := &SpyGame{}

		c := cli.NewCLI(playerStore, in, stdout, game)
		c.PlayPoker()

		assert.Equal(t, 1, len(playerStore.WinCalls))
		assert.Equal(t, "Chris", playerStore.WinCalls[0])
	})

	t.Run("it starts the game with the number of players from user input", func(t *testing.T) {
		in := strings.NewReader("7\nChris wins\n")
		game := &SpyGame{}
		store := &testhelpers.StubPlayerStore{}
		out := &bytes.Buffer{}

		c := cli.NewCLI(store, in, out, game)
		c.PlayPoker()

		assert.Equal(t, 7, game.startCalledWith)
	})
}
