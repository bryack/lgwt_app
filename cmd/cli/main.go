package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bryack/lgwt_app/adapters/cli"
	"github.com/bryack/lgwt_app/filesystem"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := filesystem.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	cli.NewCLI(store, os.Stdin).PlayPoker()

}
