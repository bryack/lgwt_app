package cli_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/bryack/lgwt_app/adapters/cli"
	"github.com/stretchr/testify/assert"
)

type SpyGame struct {
	startCalledWith  int
	finishCalledWith string
	startCalled      bool
}

func (s *SpyGame) Start(numberOfPlayers int, alertsDestination io.Writer) {
	s.startCalledWith = numberOfPlayers
	s.startCalled = true
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

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("Pies\n")
		game := &SpyGame{}

		c := cli.NewCLI(in, stdout, game)
		c.PlayPoker()

		assert.True(t, game.startCalledWith == 0)
		assert.True(t, !game.startCalled, "game should not have started")

		gotPrompt := stdout.String()
		wantPrompt := cli.PlayerPrompt + cli.BadPlayerInputErrMsg
		assert.Equal(t, wantPrompt, gotPrompt)
	})
}
