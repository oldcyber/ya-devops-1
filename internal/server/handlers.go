package server

import (
	"net/http"

	"ya-devops-1/internal/data"

	"github.com/go-chi/chi/v5"
)

var str = data.NewstoredData()

// GetRoot сервер должен отдавать HTML-страничку со списком имён и значений всех известных ему на текущий момент метрик.
func GetRoot(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	for v, k := range str.GetStoredData() {
		ik := "name: " + v + " value: " + k + "\n"
		_, err := w.Write([]byte(ik))
		if err != nil {
			return
		}
	}
}

// GetMetrics читаем данные из URL и сохраняем
func GetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	var res []string
	res = append(res, chi.URLParam(r, "type"))
	res = append(res, chi.URLParam(r, "name"))
	res = append(res, chi.URLParam(r, "value"))

	if res == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	er, an := str.AddStoredData(res)
	if !er {
		w.WriteHeader(an)
		return
	} else {
		w.WriteHeader(200)
	}
}

// GetValue должен возвращать текущее значение запрашиваемой метрики
// в текстовом виде по запросу GET
// http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
func GetValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	typeM := chi.URLParam(r, "type")
	nameM := chi.URLParam(r, "name")
	if typeM != "gauge" && typeM != "counter" {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Нет такого типа метрики"))
		if err != nil {
			return
		}
		return
	}

	res, status := str.GetStoredDataByName(typeM, nameM)

	if status != 200 {
		w.WriteHeader(status)
		return
	}

	_, err := w.Write([]byte(res))
	if err != nil {
		return
	}
}
