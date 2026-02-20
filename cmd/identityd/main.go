package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nicholasalexander/digital-identity/pkg/api"
	"github.com/nicholasalexander/digital-identity/pkg/store"
)

func main() {
	storePath := os.Getenv("IDENTITY_STORE_PATH")
	if storePath == "" {
		storePath = "./data/chains.json"
	}

	st, err := store.NewFileStore(storePath)
	if err != nil {
		log.Fatalf("initialize store: %v", err)
	}
	srv := api.NewServer(st)

	addr := os.Getenv("IDENTITY_HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("identity API listening on %s", addr)
	if err := http.ListenAndServe(addr, srv.Handler()); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
