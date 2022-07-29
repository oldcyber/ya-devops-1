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
	r.Post("/update/", server.UpdateJSONMetrics)
	r.Post("/value/", server.GetJSONMetric)
	r.Post("/update/{type}/{name}/{value}", server.UpdateMetrics)
	r.Get("/value/{type}/{name}", server.GetMetric)

	log.Error(http.ListenAndServe(":8080", r))
}
