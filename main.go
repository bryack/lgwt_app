package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bryack/lgwt_app/filesystem"
	"github.com/bryack/lgwt_app/server"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("failed to open file %q: %v", dbFileName, err)
	}

	store := filesystem.NewFileSystemPlayerStore(db)

	server := server.NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
