package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nickklius/go-short/internal/config"
	"github.com/nickklius/go-short/internal/storage"
	"io"
	"net/http"
)

type URL struct {
	URL string `json:"url"`
}

func ServiceRouter(repo storage.Repository) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", RetrieveHandler(repo))
		r.Post("/", ShortenHandler(repo))
		r.Post("/api/shorten", ShortenJsonHandler(repo))
	})

	return r
}

func ShortenJsonHandler(URLStorage storage.Repository) http.HandlerFunc {
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

		shortURL, err := storage.CreateShortURL(URLStorage, context.Background(), url.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result := struct {
			Result string `json:"result"`
		}{
			Result: config.ServiceURL + "/" + shortURL,
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

func ShortenHandler(URLStorage storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(b) > 0 {
			shortURL, err := storage.CreateShortURL(URLStorage, context.Background(), string(b))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)

			_, err = w.Write([]byte(config.ServiceURL + "/" + shortURL))
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

func RetrieveHandler(URLStorage storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "id")
		longURL, err := storage.RetrieveURL(URLStorage, context.Background(), shortURL)

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
