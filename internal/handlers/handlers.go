package handlers

import (
	"io"
	"math/rand"
	"net/http"
)

const (
	host    = "localhost"
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	port    = "8080"
	schema  = "http"
)

var sh = URLShortener{storage: map[string]string{}}

type URLShortener struct {
	storage map[string]string
}

func (u *URLShortener) generateKey(n int) string {
	b := make([]byte, n)
	for {
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		if _, ok := u.storage[string(b)]; !ok {
			break
		}
	}
	return string(b)
}

func (u *URLShortener) checkURL(url string) bool {
	if _, ok := u.storage[url]; ok {
		return true
	}
	return false
}

func (u *URLShortener) shortenURL(url string) string {
	short := u.generateKey(5)
	u.storage[short] = url
	return short
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if len(b) > 0 {
		url := sh.shortenURL(string(b))
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(schema + "://" + host + ":" + port + "/" + url))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func retrieveHandler(w http.ResponseWriter, r *http.Request) {
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
		retrieveHandler(w, r)
	case http.MethodPost:
		shortenHandler(w, r)
	default:
		http.Error(w, "Allowed only GET and POST methods", http.StatusBadRequest)
	}
}
