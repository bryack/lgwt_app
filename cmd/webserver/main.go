package main

import (
	"log"
	"net/http"

	"github.com/bryack/lgwt_app/adapters/alerter"
	"github.com/bryack/lgwt_app/adapters/server"
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

	game := game.NewGame(scheduler.BlindAlerterFunc(alerter.Alerter), store)
	server, err := server.NewPlayerServer(store, game)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
