package handlers

import (
	"context"
	"github.com/nickklius/go-short/internal/storage"
	"io"
	"net/http"
)

const (
	host   = "localhost"
	port   = "8080"
	schema = "http"
)

func ShortenHandler(URLStorage storage.Repository) http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(b) > 0 {
			shortURL, err := storage.CreateShortURL(URLStorage, ctx, string(b))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)

			_, err = w.Write([]byte(schema + "://" + host + ":" + port + "/" + shortURL))
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
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := r.URL.Path[1:]
		longURL, err := storage.RetrieveURL(URLStorage, ctx, shortURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		switch longURL != "" {
		case true:
			http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
		default:
			http.Error(w, "URL not found", http.StatusBadRequest)
		}
	}
}
