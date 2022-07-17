package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestGetMetrics(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "negative test #1",
			request: "counter/testCounter/none",
			want: want{
				code:        400,
				contentType: "text/plain",
			},
		},
		{
			name:    "negative test #2",
			request: "gauge/testGauge/none",
			want: want{
				code:        400,
				contentType: "text/plain",
			},
		},
		{
			name:    "negative test #3",
			request: "counter/",
			want: want{
				code:        404,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "positive test #1",
			request: "counter/testCounter/100",
			want: want{
				code:        200,
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StoredData = make(map[string]StoredType)
			target := "/update/" + tt.request
			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, target, nil)
			request.Header.Set("Content-Type", "text/plain")

			r := chi.NewRouter()
			r.Post("/update/{type}/{name}/{value}", GetMetrics)

			// h := http.HandlerFunc(GetMetrics)
			r.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected content type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

//func TestGetRoot(t *testing.T) {
//	type want struct {
//		code        int
//		contentType string
//		response    string
//	}
//	tests := []struct {
//		name string
//		want want
//	}{
//		{
//			name: "positive test #1",
//			want: want{
//				code:        200,
//				contentType: "text/plain",
//				response:    "Hello, stranger!",
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			request := httptest.NewRequest(http.MethodGet, "/", nil)
//			w := httptest.NewRecorder()
//			h := http.HandlerFunc(GetRoot)
//			h.ServeHTTP(w, request)
//			res := w.Result()
//			if res.StatusCode != tt.want.code {
//				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
//			}
//			defer res.Body.Close()
//			resBody, err := io.ReadAll(res.Body)
//			if err != nil {
//				t.Fatal(err)
//			}
//			if string(resBody) != tt.want.response {
//				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
//			}
//			if res.Header.Get("Content-Type") != tt.want.contentType {
//				t.Errorf("Expected content type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
//			}
//		})
//	}
//}
