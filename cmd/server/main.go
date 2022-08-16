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

func main() {
	log.Info("Starting server")
	log.Info("Checking environment variables")
	cfg := tools.NewConfig()
	if err := cfg.InitFromFlags(); err != nil {
		log.Error(err)
		return
	}
	cfg.PrintConfig()
	log.Println("loading config. Address:", cfg.Address, "Restore:", cfg.Restore, "Store interval", cfg.StoreInterval.Seconds(), "Store file", cfg.StoreFile)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(server.GzipMiddleware)
	r.Get("/", server.GetRoot)
	r.Post("/update/", server.UpdateJSONMetrics)
	r.Post("/value/", server.GetJSONMetric)
	r.Post("/update/{type}/{name}/{value}", server.UpdateMetrics)
	r.Get("/value/{type}/{name}", server.GetMetric)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		log.Error(http.ListenAndServe(cfg.Address, r))
		wg.Done()
	}()
	go func() {
		err := &cfg.Restore
		if *err {
			err := server.ReadLogFile()
			if err != nil {
				log.Error(err)
			}
		}

		if cfg.StoreFile != "" {
			err := server.WorkWithLogs(cfg)
			if err != nil {
				log.Error(err)
				return
			}
		} else {
			log.Info("Писать ничего не будем")
			return
		}
		wg.Done()
	}()
	go func() {
		<-c
		log.Info("Shutdown server")
		os.Exit(1)
	}()
	wg.Wait()
}
