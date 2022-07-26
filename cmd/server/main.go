package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"ya-devops-1/internal/server"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", server.GetRoot)
	r.Post("/update/", server.GetJSONMetrics)
	r.Post("/value/", server.GetJSONValue)
	r.Post("/update/{type}/{name}/{value}", server.GetMetrics)
	r.Get("/value/{type}/{name}", server.GetValue)

	log.Error(http.ListenAndServe(":8080", r))
}
