package app

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tim3-p/go-link-shortener/internal/configs"
	"github.com/tim3-p/go-link-shortener/internal/models"
	"github.com/tim3-p/go-link-shortener/internal/storage"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	var repository storage.Repository
	var tChan chan *models.Task
	if configs.EnvConfig.FileStoragePath == "" {
		repository = storage.NewMapRepository()
	}

	handler := NewAppHandler(repository, tChan)

	r := NewRouter(handler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	type want struct {
		response   string
		statusCode int
	}
	tests := []struct {
		name    string
		method  string
		request string
		want    want
	}{
		{
			name:   "test for POST method",
			method: http.MethodPost,
			want: want{
				statusCode: 201,
			},
			request: "/",
		},
		{
			name:   "test for GET method",
			method: http.MethodGet,
			want: want{
				statusCode: 400,
			},
			request: "/fgRth",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, tt.method, tt.request)
			defer resp.Body.Close()

			if resp.StatusCode != tt.want.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.want.statusCode, resp.StatusCode)
			}
		})
	}

}
