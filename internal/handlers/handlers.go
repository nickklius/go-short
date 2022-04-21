package handlers

import (
	"github.com/nickklius/go-short/internal/utils"
	"io"
	"net/http"
)

const (
	host   = "localhost"
	port   = "8080"
	schema = "http"
	urllen = 5
)

var sh = URLShortener{storage: map[string]string{}}

type URLShortener struct {
	storage map[string]string
}

func (u *URLShortener) checkURL(url string) bool {
	if _, ok := u.storage[url]; ok {
		return true
	}
	return false
}

func (u *URLShortener) shortenURL(url string) string {
	var short string
	for {
		short = utils.GenerateKey(urllen)
		if _, ok := u.storage[short]; !ok {
			break
		}
	}
	u.storage[short] = url
	return short
}

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if len(b) > 0 {
		url := sh.shortenURL(string(b))
		w.Header().Set("Content-Type", "plain/text")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(schema + "://" + host + ":" + port + "/" + url))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func RetrieveHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[1:]
	switch sh.checkURL(url) {
	case true:
		http.Redirect(w, r, sh.storage[url], http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func URLHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: validate URL with URLShortener method

	switch r.Method {
	case http.MethodGet:
		RetrieveHandler(w, r)
	case http.MethodPost:
		ShortenHandler(w, r)
	default:
		http.Error(w, "Allowed only GET and POST methods", http.StatusBadRequest)
	}
}
