package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/nickklius/go-short/internal/config"
	"github.com/nickklius/go-short/internal/storage"
	"io"
	"net/http"
)

type Handler struct {
	Storage storage.Repository
	Config  config.Config
}

type URL struct {
	URL string `json:"url"`
}

func New() *Handler {
	return &Handler{
		Storage: &storage.MapURLStorage{Storage: map[string]string{}},
		Config:  config.New(),
	}
}

func (h *Handler) ShortenJsonHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		url := URL{}

		err = json.Unmarshal(b, &url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if url.URL == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		shortURL, err := storage.CreateShortURL(h.Storage, context.Background(), url.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result := struct {
			Result string `json:"result"`
		}{
			Result: h.Config.BaseURL + shortURL,
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)

		b, err = json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) ShortenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(b) > 0 {
			shortURL, err := storage.CreateShortURL(h.Storage, context.Background(), string(b))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)

			_, err = w.Write([]byte(h.Config.BaseURL + shortURL))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

func (h *Handler) RetrieveHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "id")
		longURL, err := storage.RetrieveURL(h.Storage, context.Background(), shortURL)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if longURL != "" {
			http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
		} else {
			http.Error(w, "URL not found", http.StatusBadRequest)
		}
	}
}
