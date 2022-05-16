package handlers

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"github.com/nickklius/go-short/internal/config"
	"github.com/nickklius/go-short/internal/storages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var c = config.NewConfig()

func TestMain(m *testing.M) {
	c.ParseFlags()
	os.Exit(m.Run())
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestRetrieveHandler(t *testing.T) {
	s := storages.NewMemoryStorage(c)
	h := NewHandler(s, c)

	if _, err := h.storage.Create("https://yandex.ru"); err != nil {
		log.Fatal(err)
	}

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
				statusCode: http.StatusTemporaryRedirect,
			},
		},
		{
			name: "fail: wrong short url 400",
			path: "/abcde",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "fail: empty short url",
			path: "/",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRouter := chi.NewRouter()
			testRouter.Get("/{id}", h.RetrieveHandler())

			ts := httptest.NewServer(testRouter)
			defer ts.Close()

			resp, _ := testRequest(t, ts, http.MethodGet, tt.path, nil)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
		})
	}
}

func TestShortenHandler(t *testing.T) {
	s := storages.NewMemoryStorage(c)
	h := NewHandler(s, c)

	type want struct {
		statusCode    int
		lenShortenURL int
	}

	tests := []struct {
		name string
		body string
		path string
		want want
	}{
		{
			name: "success: URL shorten, response status code 201",
			body: "https://ya.ru",
			path: "/",
			want: want{
				statusCode:    http.StatusCreated,
				lenShortenURL: len(h.config.BaseURL) + 1 + h.config.KeyLength,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRouter := chi.NewRouter()
			testRouter.Post(tt.path, h.ShortenHandler())

			ts := httptest.NewServer(testRouter)
			defer ts.Close()

			resp, resultBody := testRequest(t, ts, http.MethodPost, "/", bytes.NewBuffer([]byte(tt.body)))
			defer resp.Body.Close()

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.lenShortenURL, len(resultBody))
		})
	}
}

func TestShortenJsonHandler(t *testing.T) {
	s := storages.NewMemoryStorage(c)
	h := NewHandler(s, c)

	type want struct {
		statusCode   int
		contentType  string
		responseBody string
	}

	tests := []struct {
		name string
		body string
		path string
		want want
	}{
		{
			name: "success: URL shorten json, response status code 201",
			body: "{\"url\":\"https://ya.ru/\"}",
			path: "/api/shorten",
			want: want{
				statusCode:   http.StatusCreated,
				contentType:  "application/json; charset=utf-8",
				responseBody: "{\"result\":\"http://localhost:8080/e7ut4\"}",
			},
		},
		{
			name: "fail: wrong request body for shorten json handler",
			body: "{\"_\":\"https://ya.ru/\"}",
			path: "/api/shorten",
			want: want{
				statusCode:   http.StatusBadRequest,
				contentType:  "",
				responseBody: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRouter := chi.NewRouter()
			testRouter.Post(tt.path, h.ShortenJSONHandler())

			ts := httptest.NewServer(testRouter)
			defer ts.Close()

			resp, resultBody := testRequest(t, ts, http.MethodPost, tt.path, bytes.NewBuffer([]byte(tt.body)))
			defer resp.Body.Close()

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.responseBody, resultBody)
		})
	}
}
