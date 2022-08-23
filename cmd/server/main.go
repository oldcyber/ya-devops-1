package main

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/oldcyber/ya-devops-1/internal/server"
	"github.com/oldcyber/ya-devops-1/internal/tools"
	log "github.com/sirupsen/logrus"
)

// var DBPool *pgxpool.Pool

func main() {
	log.Info("Starting server")
	log.Info("Checking environment variables")
	cfg := tools.NewConfig()
	if err := cfg.InitFromEnv(); err != nil {
		log.Error(err)
		return
	}
	if err := cfg.InitFromServerFlags(); err != nil {
		log.Error(err)
		return
	}
	cfg.PrintConfig()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(server.GzipMiddleware)
	r.Get("/ping", server.GetPing(http.HandlerFunc(server.Ping), cfg))
	r.Get("/", server.GetRoot)
	r.Post("/update/", server.CheckHash(http.HandlerFunc(server.UpdateJSONMetrics), cfg))
	// r.Post("/update/", server.CheckHash(http.HandlerFunc(server.UpdateJSONMetrics), cfg))
	r.Post("/value/", server.GetHash(http.HandlerFunc(server.GetJSONMetric), cfg))
	// r.Post("/value/", server.GetJSONMetric)
	r.Post("/update/{type}/{name}/{value}", server.UpdateMetrics)
	r.Get("/value/{type}/{name}", server.GetMetric)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		log.Error(http.ListenAndServe(cfg.GetAddress(), r))
		wg.Done()
	}()
	go func() {
		wg.Done()
		err := cfg.GetRestore()
		if err {
			err := server.ReadLogFile(cfg)
			if err != nil {
				log.Error(err)
			}
		}

		if cfg.GetStoreFile() != "" {
			err := server.WorkWithLogs(cfg)
			if err != nil {
				log.Error("Проблема с записью из main", err)
				return
			}
		} else {
			log.Info("Писать ничего не будем")
			return
		}
	}()
	go func() {
		<-c
		log.Info("Shutdown server")
		os.Exit(1)
	}()
	wg.Wait()
}
