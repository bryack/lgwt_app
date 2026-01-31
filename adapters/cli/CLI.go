package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/bryack/lgwt_app/scheduler"
	"github.com/bryack/lgwt_app/store"
)

const PlayerPrompt = "Please enter the number of players: "

type CLI struct {
	store   store.PlayerStore
	in      *bufio.Scanner
	out     io.Writer
	alerter scheduler.BlindAlerter
}

func NewCLI(store store.PlayerStore, in io.Reader, out io.Writer, alerter scheduler.BlindAlerter) *CLI {
	return &CLI{
		store:   store,
		in:      bufio.NewScanner(in),
		out:     out,
		alerter: alerter,
	}
}

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)
	cli.scheduleBlindAlerts()
	userInput := cli.readLine()
	cli.store.RecordWin(extractWinner(userInput))
}

func (cli *CLI) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + 10*time.Minute
	}
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
