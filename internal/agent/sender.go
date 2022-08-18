package agent

import (
	"bytes"
	"net/http"
	"time"

	"github.com/oldcyber/ya-devops-1/internal/data"
	log "github.com/sirupsen/logrus"
)

// Pinger проверяет доступность сервера
func Pinger(cfg config) error {
	//	Проверяем жив ли сервер
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://"+cfg.GetAddress(), nil)
	if err != nil {
		log.Println("Ошибка запроса: ", err)
	}
	// Делаем запросы с таймаутом проверки живости сервера
	for i := 0; i < 5; i++ {
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Ошибка при отправке данных в сервер: ", err)
			time.Sleep(10 * time.Second)
		} else {
			defer resp.Body.Close()
			log.Println("Сервер жив:", resp.Status)
			return nil
		}
	}
	return err
}

// sendJSONGaugeMetrics отправляет данные типа gauge в сервис метрик в формате JSON
func sendJSONGaugeMetrics(m map[string]float64, cfg config) {
	// Проверяем доступность сервера
	err := Pinger(cfg)
	if err != nil {
		log.Println("Сервер не доступен: ", err)
		return
	}
	//	Инициализируем клиента
	client := &http.Client{}
	// проходим по метрикам и отправляем их на сервер
	for key, val := range m {
		func() {
			m := data.Metrics{}
			j := m.SendGaugeMetrics(key, val)
			// retries := 3
			var resp *http.Response
			// var resp *http.Response
			req, err := http.NewRequest(http.MethodPost, "http://"+cfg.GetAddress()+"/update/", bytes.NewBuffer(j))
			if err != nil {
				log.Println("Ошибка запроса: ", err)
				// log.Println(err)
			}
			req.Header.Add("Content-Type", "application/json")
			resp, err = client.Do(req)
			if err != nil {
				log.Println("Ошибка при отправке данных в сервис метрик: ", err)
			} else {
				defer resp.Body.Close()
				log.Println("Отправлено на сервер:", string(j), "Статус-код ", resp.Status)
			}
		}()
	}
}

// sendJSONCounterMetrics отправляет данные типа counter в сервис метрик в формате JSON
func sendJSONCounterMetrics(c int64, cfg config) {
	// Проверяем доступность сервера
	err := Pinger(cfg)
	if err != nil {
		log.Println("Сервер не доступен: ", err)
		return
	}
	//	Инициализируем клиента
	client := &http.Client{}
	func() {
		m := data.Metrics{}
		j := m.SendCounterMetrics(c)
		// retries := 3
		var resp *http.Response
		req, err := http.NewRequest(http.MethodPost, "http://"+cfg.GetAddress()+"/update/", bytes.NewBuffer(j))
		if err != nil {
			log.Println(err)
		}
		req.Header.Add("Content-Type", "application/json")
		resp, err = client.Do(req)
		if err != nil {
			log.Println("Ошибка при отправке данных в сервис метрик: ", err)
		} else {
			defer resp.Body.Close()
			log.Println("Отправлено на сервер:", string(j), "Статус-код ", resp.Status)
		}
	}()
}
