package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bryack/lgwt_app/adapters/server"
	"github.com/bryack/lgwt_app/filesystem"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("failed to open file %q: %v", dbFileName, err)
	}

	store, err := filesystem.NewFileSystemPlayerStore(db)
	if err != nil {
		log.Fatalf("failed to create file system player store: %v", err)
	}

	server := server.NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
