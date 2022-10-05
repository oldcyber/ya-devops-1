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

	// Create DB
	db, _ := tools.DBConnect(cfg.GetDatabaseDSN())
	err := db.Ping()
	var myPing bool
	if err != nil {
		myPing = false
		log.Error(err)
	} else {
		myPing = true
	}
	if err != nil {
		log.Error("Ошибка соединения: ", err)
		cfg.DatabaseDSN = ""
		// return
	} else {
		err = tools.CreateTable(db)
		if err != nil {
			log.Error("Ошибка создания таблицы: ", err)
			// return
		}
	}

	defer db.Close()

	var storeTO string
	if !myPing && cfg.GetStoreFile() != "" {
		storeTO = "file"
	} else {
		storeTO = "db"
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(server.GzipMiddleware)
	r.Get("/ping", server.GetPing(http.HandlerFunc(server.Ping), db))
	r.Get("/", server.GetRoot)
	// Инкремент 2
	// StoreMetrics
	r.Post("/update/{type}/{name}/{value}", server.StoreMetrics(http.HandlerFunc(server.Plug), db, storeTO))
	// Инкремент 3
	// GetMetrics
	r.Get("/value/{type}/{name}", server.GetMetrics(http.HandlerFunc(server.Plug), db, storeTO))
	// Инкремент 4
	// StoreMetricsFromJSON
	r.Post("/update/", server.StoreMetricsFromJSON(http.HandlerFunc(server.Plug), cfg, db, storeTO))
	// GetMetricsFromJSON
	r.Post("/value/", server.GetMetricsFromJSON(http.HandlerFunc(server.Plug), cfg, db, storeTO))
	// Инкремент 12
	// MassStoreMetrics
	r.Post("/updates/", server.MassStoreMetrics(http.HandlerFunc(server.Plug), cfg, db, storeTO))

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
