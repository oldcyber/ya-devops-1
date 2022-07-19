package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"ya-devops-1/internal/server"
)

func main() {
	server.StoredData = make(map[string]server.StoredType)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", server.GetRoot)
	r.Post("/update/{type}/{name}/{value}", server.GetMetrics)
	r.Get("/value/{type}/{name}", server.GetValue)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
