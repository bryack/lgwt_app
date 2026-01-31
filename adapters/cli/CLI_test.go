package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bryack/lgwt_app/adapters/cli"
	"github.com/stretchr/testify/assert"
)

type SpyGame struct {
	startCalledWith  int
	finishCalledWith string
}

func (s *SpyGame) Start(numberOfPlayers int) {
	s.startCalledWith = numberOfPlayers
}

func (s *SpyGame) Finish(winner string) {
	s.finishCalledWith = winner
}

func TestCLI(t *testing.T) {
	t.Run("it prompts the user to enter the number of players and starts the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		game := &SpyGame{}

		c := cli.NewCLI(in, stdout, game)
		c.PlayPoker()

		assert.Equal(t, cli.PlayerPrompt, stdout.String())
		assert.Equal(t, 7, game.startCalledWith)
	})
	t.Run("it finishes the game with the winner", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("1\nChris wins\n")
		game := &SpyGame{}

		c := cli.NewCLI(in, stdout, game)
		c.PlayPoker()

		assert.Equal(t, "Chris", game.finishCalledWith)
	})
}
