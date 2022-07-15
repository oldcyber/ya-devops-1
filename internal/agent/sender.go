package agent

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"ya-devops-1/internal/tools"
)

// sendGaugeMetrics отправляет данные типа gauge в сервис метрик
func sendGaugeMetrics(m map[string]gauge) {
	//	Инициализируем клиента
	client := &http.Client{}
	// проходим по метрикам и отправляем их на сервер
	for key, val := range m {
		func() {
			req, err := http.NewRequest(http.MethodPost, tools.SetURL(key, fmt.Sprintf("%v", val), "gauge"), nil)
			if err != nil {
				log.Fatalln(err)
			}
			req.Header.Add("Content-Type", "text/plain")
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			log.Println("Статус-код ", resp.Status)
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Fatalln(err)
				}
			}(resp.Body)
		}()
	}
}

// sendCounterMetrics отправляет данные типа counter в сервис метрик
func sendCounterMetrics(c counter) {
	//	Инициализируем клиента
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, tools.SetURL("PollCount", strconv.FormatInt(int64(c), 10), "counter"), nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Content-Type", "text/plain")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("Статус-код ", resp.Status)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)
}
