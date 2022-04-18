package app

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

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

func URLHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: make post method handler

	if r.Method == http.MethodGet {
		url := r.URL.Path[1:]
		switch sh.checkURL(url) {
		case true:
			http.Redirect(w, r, sh.storage[url], http.StatusMovedPermanently)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	if r.Method == http.MethodPost {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if len(b) > 0 {
			url := sh.shortenURL(string(b))
			fmt.Println(sh.storage)
			w.Write([]byte(url))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
