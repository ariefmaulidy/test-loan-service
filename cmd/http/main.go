package main

import (
	"log"
	"net/http"

	"github.com/ariefmaulidy/test-loan-service/server"
	httpserver "github.com/ariefmaulidy/test-loan-service/server/http"
)

func main() {
	deps := server.InitDependencies()

	srv := httpserver.New(deps)

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, srv); err != nil {
		log.Fatal(err)
	}
}
