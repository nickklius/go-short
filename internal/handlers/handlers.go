package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nickklius/go-short/internal/config"
	"github.com/nickklius/go-short/internal/storages"
)

type Handler struct {
	storage storages.Repository
	config  config.Config
}

type URL struct {
	URL string `json:"url"`
}

func NewHandler(s storages.Repository, c config.Config) *Handler {
	return &Handler{
		storage: s,
		config:  c,
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
			shortURL, err := storages.CreateShortURL(h.storage, string(b))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)

			_, err = w.Write([]byte(h.config.BaseURL + "/" + shortURL))
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

func (h *Handler) ShortenJSONHandler() http.HandlerFunc {
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

		shortURL, err := storages.CreateShortURL(h.storage, url.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result := struct {
			Result string `json:"result"`
		}{
			Result: h.config.BaseURL + "/" + shortURL,
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

func (h *Handler) RetrieveHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "id")
		longURL, err := storages.RetrieveURL(h.storage, shortURL)

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
