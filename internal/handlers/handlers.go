package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/nickklius/go-short/internal/config"
	"github.com/nickklius/go-short/internal/middleware"
	"github.com/nickklius/go-short/internal/model"
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

	userID, err := middleware.GetCurrentUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortURL, err := h.prepareShortening(r.Context(), string(b), userID)

	var e *storages.InsertURLUniqError

	if errors.As(err, &e) {
		w.WriteHeader(http.StatusConflict)
		_, err = w.Write([]byte(h.config.BaseURL + "/" + e.ShortURL))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

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
	type result struct {
		Result string `json:"result"`
	}

	u := struct {
		URL string `json:"url"`
	}{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, err := middleware.GetCurrentUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respHTTPStatus := http.StatusCreated

	shortURL, err := h.prepareShortening(r.Context(), u.URL, userID)
	if err != nil {
		var e *storages.InsertURLUniqError

		if errors.As(err, &e) {
			shortURL = e.ShortURL
			respHTTPStatus = http.StatusConflict
		} else {
			http.Error(w, err.Error(), errToStatus(err))
			return
		}
	}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(result{Result: h.config.BaseURL + "/" + shortURL})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(respHTTPStatus)
	_, err = w.Write(buf.Bytes())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ShortenJSONBatchHandler(w http.ResponseWriter, r *http.Request) {
	var urls []model.URLBatchRequest

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, err := middleware.GetCurrentUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]model.URLBatchResponse, 0, len(urls))

	for _, u := range urls {
		shortURL, err := h.prepareShortening(r.Context(), u.OriginalURL, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := model.URLBatchResponse{
			CorrelationID: u.CorrelationID,
			ShortURL:      h.config.BaseURL + "/" + shortURL,
		}

		result = append(result, response)
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
	longURL, err := h.storage.Read(r.Context(), shortURL)

	if err == nil {
		http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
		return
	}

	status := errToStatus(err)
	http.Error(w, err.Error(), status)
}

func (h *Handler) RetrieveUserURLs(w http.ResponseWriter, r *http.Request) {
	type result struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	var response []result

	userID, err := middleware.GetCurrentUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	urls, err := h.storage.GetAllByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		http.Error(w, "", http.StatusNoContent)
		return
	}

	for short, long := range urls {
		response = append(response, result{
			ShortURL:    h.config.BaseURL + "/" + short,
			OriginalURL: long})
	}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf.Bytes())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PingDB(w http.ResponseWriter, _ *http.Request) {
	err := h.storage.Ping()
	if err != nil {
		http.Error(w, err.Error(), errToStatus(err))
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) prepareShortening(ctx context.Context, longURL, userID string) (string, error) {
	var shortURL string

	_, err := url.ParseRequestURI(longURL)
	if err != nil {
		return shortURL, ErrWrongURLFormat
	}

	for i := 0; i < h.config.ShortenerCapacity; i++ {
		shortURL = utils.GenerateKey(h.config.Letters, h.config.KeyLength)
		err = h.storage.Create(ctx, shortURL, longURL, userID)
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
	case storages.ErrMethodNotImplemented:
		return http.StatusNotImplemented
	default:
		return http.StatusInternalServerError
	}
}
