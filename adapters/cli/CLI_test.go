package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bryack/lgwt_app/adapters/cli"
	"github.com/bryack/lgwt_app/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestCLI(t *testing.T) {
	t.Run("it prompts the user to enter the number of players and starts the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		game := &testhelpers.SpyGame{}

		c := cli.NewCLI(in, stdout, game)
		c.PlayPoker()

		assert.Equal(t, cli.PlayerPrompt, stdout.String())
		assert.Equal(t, 7, game.StartCalledWith)
	})
	t.Run("it finishes the game with the winner", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("1\nChris wins\n")
		game := &testhelpers.SpyGame{}

		c := cli.NewCLI(in, stdout, game)
		c.PlayPoker()

		assert.Equal(t, "Chris", game.FinishCalledWith)
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("Pies\n")
		game := &testhelpers.SpyGame{}

		c := cli.NewCLI(in, stdout, game)
		c.PlayPoker()

		assert.True(t, game.StartCalledWith == 0)
		assert.True(t, !game.StartCalled, "game should not have started")

		gotPrompt := stdout.String()
		wantPrompt := cli.PlayerPrompt + cli.BadPlayerInputErrMsg
		assert.Equal(t, wantPrompt, gotPrompt)
	})
}
