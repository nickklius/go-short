package service

import (
	"context"
	"net/http"
	"time"

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
	done    chan struct{}
}

func NewService(ctx context.Context) (*Service, error) {
	var s storages.Repository

	c, err := config.NewConfig()
	done := make(chan struct{})

	if err != nil {
		return nil, err
	}

	switch {
	case c.DatabaseDSN != "":
		s, err = storages.NewDatabaseStorage(ctx, c, done)
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
		done:    done,
	}, nil
}

func (s *Service) Start(ctx context.Context, errCh chan error) {
	h := handlers.NewHandler(s.Storage, s.Conf)
	srv := &http.Server{Addr: s.Conf.ServerAddress, Handler: s.Router(h)}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	<-ctx.Done()

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		cancel()
	}()

	srv.SetKeepAlivesEnabled(false)

	if err := srv.Shutdown(ctxTimeout); err != nil {
		errCh <- err
	}

	<-s.done
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
