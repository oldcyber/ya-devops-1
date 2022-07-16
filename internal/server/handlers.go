package server

import (
	"log"
	"net/http"

	"ya-devops-1/internal/tools"
)

// GetRoot - обработчик запроса на главную страницу
//func GetRoot(w http.ResponseWriter, _ *http.Request) {
//	w.Header().Set("Content-Type", "text/plain")
//
//	_, err := w.Write([]byte("Hello, stranger!"))
//	if err != nil {
//		return
//	}
//}

// GetMetrics читаем данные из URL и сохраняем
func GetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	res := tools.GetURL(r.URL.Path)
	er, an := storeData(res)
	if !er {
		w.WriteHeader(an)
		return
	} else {
		w.WriteHeader(200)
	}

	for k, v := range StoredData {
		log.Println("key", k, "value", v)
	}
}
