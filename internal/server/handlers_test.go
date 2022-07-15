package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMetrics(t *testing.T) {
	type want struct {
		code        int
		contentType string
		response    string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "positive test #1",
			request: "gauge/Alloc/12345",
			want: want{
				code:        200,
				contentType: "text/plain",
				response:    "Я всё записал",
			},
		},
		{
			name:    "negative test #1",
			request: "test/Alloc/",
			want: want{
				code:        400,
				contentType: "text/plain",
				response:    "Что-то пошло не так",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := "/update/" + tt.request
			request := httptest.NewRequest(http.MethodPost, target, nil)
			request.Header.Set("Content-Type", "text/plain")
			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetMetrics)
			h.ServeHTTP(w, request)
			res := w.Result()
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					t.Fatal(err)
				}
			}(res.Body)
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected content type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func TestGetRoot(t *testing.T) {
	type want struct {
		code        int
		contentType string
		response    string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        200,
				contentType: "text/plain",
				response:    "Hello, stranger!",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetRoot)
			h.ServeHTTP(w, request)
			res := w.Result()
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					t.Fatal(err)
				}
			}(res.Body)
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected content type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
