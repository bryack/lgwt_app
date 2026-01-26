package main

import (
	"log"
	"net/http"

	"github.com/bryack/lgwt_app/server"
)

func main() {
	server := server.NewPlayerServer(server.NewInMemoryPlayerStore())
	log.Fatal(http.ListenAndServe(":5000", server))
}
