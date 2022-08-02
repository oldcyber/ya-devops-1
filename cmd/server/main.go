package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"ya-devops-1/internal/tools"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"ya-devops-1/internal/server"
)

func main() {
	checkEnv := func(key string) bool {
		_, ok := os.LookupEnv(key)
		if !ok {
			return false
		} else {
			return true
		}
	}
	Address := flag.String("a", "", "address")
	Restore := flag.Bool("r", true, "restore")
	StoreInterval := flag.Duration("i", 0, "store interval")
	StoreFile := flag.String("f", "", "store file")
	flag.Parse()
	log.Info("StoreInterval: ", *StoreInterval, " tools.Conf.StoreInterval: ", tools.Conf.StoreInterval)
	if !checkEnv("ADDRESS") && *Address != "" {
		tools.Conf.Address = *Address
	}
	if !checkEnv("RESTORE") && *Restore != tools.Conf.Restore {
		tools.Conf.Restore = *Restore
	}
	if !checkEnv("STORE_INTERVAL") && *StoreInterval != 0 {
		tools.Conf.StoreInterval = *StoreInterval
	}
	if !checkEnv("STORE_FILE") && *StoreFile != "" {
		tools.Conf.StoreFile = *StoreFile
	}
	log.Println("loading config. Address:", tools.Conf.Address, "Restore:", tools.Conf.Restore, "Store interval", tools.Conf.StoreInterval.Seconds(), "Store file", tools.Conf.StoreFile)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
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
		log.Error(http.ListenAndServe(tools.Conf.Address, r))
		wg.Done()
	}()
	go func() {
		err := tools.Conf.Restore
		if err {
			err := server.ReadLogFile()
			if err != nil {
				log.Error(err)
			}
		}

		if tools.Conf.StoreFile != "" {
			err := server.WorkWithLogs()
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
		//err := tools.CloseFile(server.FO)
		//if err != nil {
		//	log.Error(err)
		//}
		log.Info("Shutdown server")
		os.Exit(1)
	}()
	wg.Wait()
}
