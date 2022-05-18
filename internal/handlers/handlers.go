package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nickklius/go-short/internal/config"
	"github.com/nickklius/go-short/internal/storages"
	"github.com/nickklius/go-short/internal/utils"
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

func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(b) > 0 {
		shortURL := utils.GenerateKey(h.config.Letters, h.config.KeyLength)
		err = h.storage.Create(shortURL, string(b))
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
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

func (h *Handler) ShortenJSONHandler(w http.ResponseWriter, r *http.Request) {
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

	shortURL := utils.GenerateKey(h.config.Letters, h.config.KeyLength)
	err = h.storage.Create(shortURL, url.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
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

func (h *Handler) RetrieveHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "id")
	longURL, err := h.storage.Read(shortURL)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}
