package service

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nickklius/go-short/internal/handlers"
	"log"
	"net/http"
)

type service struct{}

func New() *service {
	return &service{}
}

func (s *service) Start() {
	h := handlers.New()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", h.RetrieveHandler())
		r.Post("/", h.ShortenHandler())
		r.Post("/api/shorten", h.ShortenJsonHandler())
	})

	log.Fatal(http.ListenAndServe(h.Config.ServerAddress, r))
}
