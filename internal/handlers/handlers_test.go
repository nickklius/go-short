package handlers

import (
	"bytes"
	"github.com/nickklius/go-short/internal/storage"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testStorage storage.Repository = &storage.MapURLStorage{Storage: map[string]string{
	"5fbbd": "https://yandex.ru",
}}

func TestRetrieveHandler(t *testing.T) {
	type want struct {
		statusCode int
	}

	tests := []struct {
		name string
		path string
		want want
	}{
		{
			name: "success: correct redirect 307",
			path: "/5fbbd",
			want: want{
				statusCode: 307,
			},
		},
		{
			name: "fail: wrong short url 400",
			path: "/abcde",
			want: want{
				statusCode: 400,
			},
		},
		{
			name: "fail: empty short url",
			path: "/",
			want: want{
				statusCode: 400,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			h := RetrieveHandler(testStorage)

			h.ServeHTTP(w, request)
			result := w.Result()

			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestShortenHandler(t *testing.T) {
	type want struct {
		statusCode    int
		lenShortenURL int
	}

	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "success: URL shorten, response status code 201",
			body: "https://ya.ru",
			want: want{
				statusCode:    201,
				lenShortenURL: len("http://localhost:8080") + 6,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(tt.body)))
			w := httptest.NewRecorder()
			h := ShortenHandler(testStorage)
			h.ServeHTTP(w, request)
			result := w.Result()

			resultBody, err := io.ReadAll(result.Body)
			defer result.Body.Close()
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.lenShortenURL, len(string(resultBody)))
		})
	}
}
