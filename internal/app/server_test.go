package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCommonHandler(t *testing.T) {
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
			name:   "test for DELETE method",
			method: http.MethodDelete,
			want: want{
				statusCode: 400,
			},
			request: "/",
		},
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
				statusCode: 307,
			},
			request: "/fgRth",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(CommonHandler)
			h.ServeHTTP(w, request)
			res := w.Result()

			if res.StatusCode != tt.want.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.want.statusCode, w.Code)
			}
		})
	}
}
