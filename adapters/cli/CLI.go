package cli

import (
	"io"

	"github.com/bryack/lgwt_app/store"
)

type CLI struct {
	store store.PlayerStore
	in    io.Reader
}

func (cli *CLI) PlayPoker() {
	cli.store.RecordWin("Chris")
}
