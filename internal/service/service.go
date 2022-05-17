package service

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nickklius/go-short/internal/config"
	"github.com/nickklius/go-short/internal/handlers"
	"github.com/nickklius/go-short/internal/storages"
)

type Service struct {
	Storage storages.Repository
	Conf    config.Config
}

func NewService() *Service {
	var s storages.Repository
	c := config.NewConfig()
	c.ParseFlags()

	if c.FileStoragePath != "" {
		s, _ = storages.NewLocalStorage(c)
	} else {
		s = storages.NewMemoryStorage(c)
	}
	return &Service{
		Storage: s,
		Conf:    c,
	}
}

func (s *Service) Start() {
	h := handlers.NewHandler(s.Storage, s.Conf)

	log.Fatal(http.ListenAndServe(s.Conf.ServerAddress, s.Router(h)))
}

func (s *Service) Router(h *handlers.Handler) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(handlers.GzipDecompressor)
	r.Use(handlers.GzipCompressor)

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", h.RetrieveHandler())
		r.Post("/", h.ShortenHandler())
		r.Post("/api/shorten", h.ShortenJSONHandler())
	})

	return r
}
