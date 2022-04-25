package handlers

import (
	"bytes"
	"github.com/nickklius/go-short/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testStorage storage.Repository = &storage.MapURLStorage{Storage: map[string]string{
		"5fbbd": "https://yandex.ru",
	}}
	testRouter = ServiceRouter(testStorage)
)

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
				statusCode: 405,
			},
		},
	}

	ts := httptest.NewServer(testRouter)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, http.MethodGet, tt.path, nil)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
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

	ts := httptest.NewServer(testRouter)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, resultBody := testRequest(t, ts, http.MethodPost, "/", bytes.NewBuffer([]byte(tt.body)))
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.lenShortenURL, len(resultBody))
		})
	}
}
