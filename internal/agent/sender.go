package agent

import (
	"bytes"
	"context"
	"net/http"

	"github.com/oldcyber/ya-devops-1/internal/storage"
	log "github.com/sirupsen/logrus"
)

// sendJSONGaugeMetrics отправляет данные типа gauge в сервис метрик в формате JSON
func sendJSONGaugeMetrics(m map[string]float64, cfg config) {
	//	Инициализируем клиента
	client := &http.Client{}
	// проходим по метрикам и отправляем их на сервер
	for k, val := range m {
		func() {
			mm := storage.Metrics{}
			j := mm.MarshalGaugeMetrics(k, cfg.GetKey(), val)
			var resp *http.Response
			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://"+cfg.GetAddress()+"/update/", bytes.NewBuffer(j))
			if err != nil {
				log.Error("Ошибка запроса: ", err)
			}
			req.Header.Add("Content-Type", "application/json")
			resp, err = client.Do(req)
			if err != nil {
				log.Error("Ошибка при отправке данных в сервис метрик: ", err)
			} else {
				defer resp.Body.Close()
			}
		}()
	}
	log.Info("Отправлено на сервер: ", len(m), " метрик")
}

// sendJSONCounterMetrics отправляет данные типа counter в сервис метрик в формате JSON
func sendJSONCounterMetrics(c int64, cfg config) {
	//	Инициализируем клиента
	client := &http.Client{}
	func() {
		m := storage.Metrics{}
		j := m.MarshalCounterMetrics(c, cfg.GetKey())
		// retries := 3
		var resp *http.Response
		ctx := context.Background()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://"+cfg.GetAddress()+"/update/", bytes.NewBuffer(j))
		if err != nil {
			log.Error(err)
		}
		req.Header.Add("Content-Type", "application/json")
		resp, err = client.Do(req)
		if err != nil {
			log.Error("Ошибка при отправке данных в сервис метрик: ", err)
		} else {
			defer resp.Body.Close()
		}
	}()
	log.Info("Отправлено на сервер: ", c, " метрик")
}

func sendBulkJSONMetrics(myMap map[string]float64, cfg config) {
	//	Инициализируем клиента
	client := &http.Client{}

	func() {
		m := storage.Metrics{}
		j := m.SendBulkMetrics(myMap)

		var resp *http.Response
		ctx := context.Background()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://"+cfg.GetAddress()+"/updates/", bytes.NewBuffer(j))
		if err != nil {
			log.Error(err)
		}
		req.Header.Add("Content-Type", "application/json")
		resp, err = client.Do(req)
		if err != nil {
			log.Error("Ошибка при отправке данных в сервис метрик: ", err)
		} else {
			defer resp.Body.Close()
		}
	}()
	log.Info("Отправлено на сервер: ", len(myMap), " метрик")
}
