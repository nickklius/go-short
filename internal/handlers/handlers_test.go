package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestURLHandler(t *testing.T) {
	url := "https://yandex.ru"

	type want struct {
		statusCodeGet   int
		statusCodePost  int
		contentTypePost string
		url             string
	}

	tests := []struct {
		name string
		want want
		body io.Reader
	}{
		{
			name: "success GET and POST endpoint test",
			want: want{
				statusCodeGet:   307,
				statusCodePost:  201,
				contentTypePost: "plain/text",
				url:             url,
			},
			body: strings.NewReader(url),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rqPost := httptest.NewRequest(http.MethodPost, "/", tt.body)
			rqPost.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			wh := http.HandlerFunc(ShortenHandler)
			wh.ServeHTTP(w, rqPost)
			rsPost := w.Result()

			assert.Equal(t, tt.want.statusCodePost, rsPost.StatusCode)
			assert.Equal(t, tt.want.contentTypePost, rsPost.Header.Get("Content-Type"))

			shortenResult, err := ioutil.ReadAll(rsPost.Body)
			require.NoError(t, err)
			err = rsPost.Body.Close()
			require.NoError(t, err)

			rqGet := httptest.NewRequest(http.MethodGet, string(shortenResult), nil)
			g := httptest.NewRecorder()
			gh := http.HandlerFunc(RetrieveHandler)
			gh.ServeHTTP(g, rqGet)
			rsGet := g.Result()

			assert.Equal(t, tt.want.statusCodeGet, rsGet.StatusCode)
			require.NotNil(t, g.Header().Get("Location"))
			assert.Equal(t, tt.want.url, g.Header().Get("Location"))

		})
	}
}
