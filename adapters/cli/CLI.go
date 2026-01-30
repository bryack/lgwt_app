package cli

import (
	"bufio"
	"io"
	"strings"
	"time"

	"github.com/bryack/lgwt_app/store"
)

type CLI struct {
	store   store.PlayerStore
	in      *bufio.Scanner
	alerter BlindAlerter
}

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

func NewCLI(store store.PlayerStore, in io.Reader, alerter BlindAlerter) *CLI {
	return &CLI{
		store:   store,
		in:      bufio.NewScanner(in),
		alerter: alerter,
	}
}

func (cli *CLI) PlayPoker() {
	cli.alerter.ScheduleAlertAt(5*time.Second, 100)
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
