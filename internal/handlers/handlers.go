package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/nickklius/go-short/internal/config"
	"github.com/nickklius/go-short/internal/storages"
	"github.com/nickklius/go-short/internal/utils"
)

var (
	ErrWrongURLFormat = errors.New("wrong format")
	ErrOverCapacity   = errors.New("shortener capacity is over")
)

type Handler struct {
	storage storages.Repository
	config  config.Config
}

func NewHandler(s storages.Repository, c config.Config) *Handler {
	return &Handler{
		storage: s,
		config:  c,
	}
}

func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortURL, err := h.prepareShortening(string(b))
	if err != nil {
		http.Error(w, err.Error(), errToStatus(err))
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(h.config.BaseURL + "/" + shortURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ShortenJSONHandler(w http.ResponseWriter, r *http.Request) {
	u := struct {
		URL string `json:"url"`
	}{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortURL, err := h.prepareShortening(u.URL)
	if err != nil {
		http.Error(w, err.Error(), errToStatus(err))
		return
	}

	result := struct {
		Result string `json:"result"`
	}{
		Result: h.config.BaseURL + "/" + shortURL,
	}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(buf.Bytes())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RetrieveHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "id")
	longURL, err := h.storage.Read(shortURL)

	if err == nil {
		http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
		return
	}

	status := errToStatus(err)
	http.Error(w, err.Error(), status)
}

func (h *Handler) prepareShortening(u string) (string, error) {
	var shortURL string

	_, err := url.ParseRequestURI(u)
	if err != nil {
		return shortURL, ErrWrongURLFormat
	}

	for i := 0; i < h.config.ShortenerCapacity; i++ {
		shortURL = utils.GenerateKey(h.config.Letters, h.config.KeyLength)
		err = h.storage.Create(shortURL, u)
		if err != storages.ErrAlreadyExists {
			break
		}
	}

	if err == storages.ErrAlreadyExists {
		return shortURL, ErrOverCapacity
	}

	if err != nil {
		return shortURL, err
	}

	return shortURL, nil
}

func errToStatus(err error) int {
	switch err {
	case ErrWrongURLFormat:
		return http.StatusBadRequest
	case ErrOverCapacity:
		return http.StatusInternalServerError
	case storages.ErrNotFound:
		return http.StatusNotFound
	case storages.ErrAlreadyExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
