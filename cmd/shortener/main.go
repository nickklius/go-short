package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nickklius/go-short/internal/handlers"
	"github.com/nickklius/go-short/internal/storage"
	"log"
	"net/http"
)

func main() {
	var URLStorage storage.Repository = &storage.MapURLStorage{Storage: map[string]string{}}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.ShortenHandler(URLStorage))
		r.Get("/", handlers.RetrieveHandler(URLStorage))
		r.Get("/{id}", handlers.RetrieveHandler(URLStorage))
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
