package main

import (
	"log"
	"net/http"

	"github.com/bryack/lgwt_app/adapters/server"
	"github.com/bryack/lgwt_app/filesystem"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := filesystem.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	server, err := server.NewPlayerServer(store)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
