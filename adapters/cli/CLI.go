package cli

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/bryack/lgwt_app/store"
)

const PlayerPrompt = "Please enter the number of players: "

type Game interface {
	Start(numberOfPlayers int)
}

type CLI struct {
	store store.PlayerStore
	in    *bufio.Scanner
	out   io.Writer
	game  Game
}

func NewCLI(store store.PlayerStore, in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		store: store,
		in:    bufio.NewScanner(in),
		out:   out,
		game:  game,
	}
}

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayers, _ := strconv.Atoi(cli.readLine())

	cli.game.Start(numberOfPlayers)
	userInput := cli.readLine()
	cli.store.RecordWin(extractWinner(userInput))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
