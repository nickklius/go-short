package service

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nickklius/go-short/internal/config"
	"github.com/nickklius/go-short/internal/handlers"
	mw "github.com/nickklius/go-short/internal/middleware"
	"github.com/nickklius/go-short/internal/storages"
)

type Service struct {
	Storage storages.Repository
	Conf    config.Config
}

func NewService(ctx context.Context) (*Service, error) {
	var s storages.Repository

	c, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	switch {
	case c.DatabaseDSN != "":
		s, err = storages.NewDatabaseStorage(ctx, c.DatabaseDSN)
		if err != nil {
			return nil, storages.ErrDBConnNotEstablished
		}
	case c.FileStoragePath != "":
		s, err = storages.NewLocalStorage(ctx, c.FileStoragePath)
		if err != nil {
			return nil, storages.ErrLocalStorageNotCreated
		}
	default:
		s = storages.NewMemoryStorage()
	}

	return &Service{
		Storage: s,
		Conf:    c,
	}, nil
}

func (s *Service) Start() error {
	h := handlers.NewHandler(s.Storage, s.Conf)

	err := http.ListenAndServe(s.Conf.ServerAddress, s.Router(h))
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Router(h *handlers.Handler) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(mw.GzipDecompressor)
	r.Use(mw.GzipCompressor)
	r.Use(mw.UserID)

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", h.RetrieveHandler)
		r.Get("/api/user/urls", h.RetrieveUserURLs)
		r.Get("/ping", h.PingDB)
		r.Post("/", h.ShortenHandler)
		r.Post("/api/shorten", h.ShortenJSONHandler)
		r.Post("/api/shorten/batch", h.ShortenJSONBatchHandler)
		r.Delete("/api/user/urls", h.DeleteURLs)
	})

	return r
}
