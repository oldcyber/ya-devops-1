package agent

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"ya-devops-1/internal/data"
)

func Pinger() error {
	//	Проверяем жив ли сервер
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080", nil)
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

// sendGaugeMetrics отправляет данные типа gauge в сервис метрик
//func sendGaugeMetrics(m map[string]float64) {
//	//	Инициализируем клиента
//	client := &http.Client{}
//	// проходим по метрикам и отправляем их на сервер
//	for key, val := range m {
//		func() {
//			req, err := http.NewRequest(http.MethodPost, tools.SetURL(key, strconv.FormatFloat(val, 'f', -1, 64), "gauge"), nil)
//			if err != nil {
//				log.Fatalln(err)
//			}
//			req.Header.Add("Content-Type", "text/plain")
//			resp, err := client.Do(req)
//			if err != nil {
//				fmt.Println(err)
//				os.Exit(1)
//			}
//			defer resp.Body.Close()
//			log.Println("Статус-код ", resp.Status)
//		}()
//	}
//}

// sendJSONGaugeMetrics отправляет данные типа gauge в сервис метрик в формате JSON
func sendJSONGaugeMetrics(m map[string]float64) {
	// Проверяем доступность сервера
	err := Pinger()
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
			req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/update/", bytes.NewBuffer(j))
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
				log.Println("Отправлено на сервер:", m.ID, m.MType, "Статус-код ", resp.Status)
			}
			//for retries > 0 {
			//	resp, err = client.Do(req)
			//	if err != nil {
			//		log.Println("Ошибка при отправке данных в сервис метрик: ", err)
			//		retries -= 1
			//		log.Println("Повторная попытка отправки данных в сервис метрик", retries)
			//		time.Sleep(10 * time.Second)
			//	} else {
			//		defer resp.Body.Close()
			//		log.Println("Отправлено на сервер:", m.ID, m.MType, "Статус-код ", resp.Status)
			//	}
			//}
		}()
	}
}

// sendCounterMetrics отправляет данные типа counter в сервис метрик
//func sendCounterMetrics(c int64) {
//	//	Инициализируем клиента
//	client := &http.Client{}
//	func() {
//		req, err := http.NewRequest(http.MethodPost, tools.SetURL("PollCount", strconv.FormatInt(int64(c), 10), "counter"), nil)
//		if err != nil {
//			log.Fatalln(err)
//		}
//		req.Header.Add("Content-Type", "text/plain")
//		resp, err := client.Do(req)
//		if err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//		defer resp.Body.Close()
//		log.Println("Статус-код ", resp.Status)
//	}()
//}

// sendJSONCounterMetrics отправляет данные типа counter в сервис метрик в формате JSON
func sendJSONCounterMetrics(c int64) {
	// Проверяем доступность сервера
	err := Pinger()
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
		req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/update/", bytes.NewBuffer(j))
		if err != nil {
			log.Println(err)
		}
		req.Header.Add("Content-Type", "application/json")
		resp, err = client.Do(req)
		if err != nil {
			log.Println("Ошибка при отправке данных в сервис метрик: ", err)
		} else {
			defer resp.Body.Close()
			log.Println("Отправлено на сервер:", m.ID, m.MType, "Статус-код ", resp.Status)
		}
		//for retries > 0 {
		//	resp, err = client.Do(req)
		//	if err != nil {
		//		log.Println("Ошибка при отправке данных в сервис метрик: ", err)
		//		retries -= 1
		//		log.Println("Повторная попытка отправки данных в сервис метрик", retries)
		//		time.Sleep(10 * time.Second)
		//	} else {
		//		defer resp.Body.Close()
		//		log.Println("Отправлено на сервер:", m.ID, m.MType, "Статус-код ", resp.Status)
		//	}
		//}
	}()
}
