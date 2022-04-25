package main

import (
	"github.com/nickklius/go-short/internal/handlers"
	"github.com/nickklius/go-short/internal/storage"
	"log"
	"net/http"
)

func main() {
	var URLStorage storage.Repository = &storage.MapURLStorage{Storage: map[string]string{}}

	r := handlers.ServiceRouter(URLStorage)
	log.Fatal(http.ListenAndServe(":8080", r))
}
