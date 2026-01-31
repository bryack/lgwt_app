package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bryack/lgwt_app/adapters/alerter"
	"github.com/bryack/lgwt_app/adapters/cli"
	"github.com/bryack/lgwt_app/filesystem"
	"github.com/bryack/lgwt_app/game"
	"github.com/bryack/lgwt_app/scheduler"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := filesystem.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	g := game.NewGame(scheduler.BlindAlerterFunc(alerter.StdOutAlerter), store)

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	cli.NewCLI(os.Stdin, os.Stdout, g).PlayPoker()

}
