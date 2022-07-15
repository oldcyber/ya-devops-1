package server

import (
	"log"
	"net/http"
	"strconv"

	"ya-devops-1/internal/tools"
)

var (
	storedData map[string]gauge
	counter    = 0
)

type gauge float64

// GetRoot - обработчик запроса на главную страницу
func GetRoot(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	_, err := w.Write([]byte("Hello, stranger!"))
	if err != nil {
		return
	}
}

// GetMetrics читаем данные из URL и сохраняем
func GetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	res := tools.GetURL(r.URL.Path)
	storedData = make(map[string]gauge)
	switch res[0] {
	case "gauge":
		// Хранить последние полученные данные вида [название][значение] ([string]gauge)
		_, err := storeData(res)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("Я всё записал"))
			if err != nil {
				return
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Я всё записал"))
		if err != nil {
			return
		}

	case "counter":
		// Прибавлять значение счётчика к предыдущему
		c, err := strconv.Atoi(res[2])
		if err != nil {
			log.Println(err)
		}
		counter = counter + c
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Я всё записал"))
		if err != nil {
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Что-то пошло не так"))
		if err != nil {
			return
		}
	}
	for k, v := range storedData {
		log.Println("key", k, "value", v)
	}
}
