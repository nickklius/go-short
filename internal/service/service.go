package service

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nickklius/go-short/internal/config"
	"github.com/nickklius/go-short/internal/handlers"
	"github.com/nickklius/go-short/internal/storages"
	"log"
	"net/http"
)

type Service struct {
	storage storages.Repository
	conf    config.Config
}

func New() *Service {
	var s storages.Repository
	var c config.Config

	c = config.New()

	if c.FileStoragePath != "" {
		s, _ = storages.NewLocalStorage(c.FileStoragePath)
	} else {
		s = storages.NewMemoryStorage()
	}
	return &Service{
		storage: s,
		conf:    c,
	}
}

func (s *Service) Start() {
	h := handlers.NewHandler(s.storage)

	log.Fatal(http.ListenAndServe(s.conf.ServerAddress, s.Router(h)))
}

func (s *Service) Router(h *handlers.Handler) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", h.RetrieveHandler())
		r.Post("/", h.ShortenHandler())
		r.Post("/api/shorten", h.ShortenJSONHandler())
	})

	return r
}
