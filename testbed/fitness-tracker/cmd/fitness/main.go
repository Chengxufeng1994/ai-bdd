// Command fitness serves the fitness tracker HTTP API.
//
// This file is the composition root: it is the only place allowed to know about
// every layer at once. Keep it to wiring — any decision made here is a decision
// that cannot be tested without starting a process.
package main

import (
	"log"
	"net/http"
	"os"

	apihttp "fitness-tracker/internal/interfaces/http"
)

// version is overridden at build time:
//
//	go build -ldflags "-X main.version=$(git describe --tags)" ./cmd/fitness
//
// The default is deliberately not a plausible version number: seeing "dev" in a
// deployed environment should look wrong immediately.
var version = "dev"

func main() {
	addr := os.Getenv("FITNESS_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	router := apihttp.NewRouter(apihttp.NewServer(version))

	log.Printf("fitness %s listening on %s", version, addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
