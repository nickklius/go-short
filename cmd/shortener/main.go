package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nickklius/go-short/internal/handlers"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{id}", handlers.RetrieveHandler)
	r.Post("/", handlers.ShortenHandler)

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
