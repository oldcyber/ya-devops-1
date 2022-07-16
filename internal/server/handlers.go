package server

import (
	"log"
	"net/http"
	"strconv"

	"ya-devops-1/internal/tools"
)

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
	er, an := storeData(res)
	if !er {
		w.WriteHeader(an)
		return
	}

	switch res[0] {
	case "gauge":
		// Хранить последние полученные данные вида [название][значение] ([string]gauge)
		er, an := storeData(res)
		if !er {
			w.WriteHeader(an)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Я всё записал"))
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
		w.WriteHeader(http.StatusNotImplemented)
		_, err := w.Write([]byte("Что-то пошло не так"))
		if err != nil {
			return
		}
	}
	for k, v := range StoredData {
		log.Println("key", k, "value", v)
	}
}
