package server

import (
	"net/http"
	"strconv"

	"ya-devops-1/internal/tools"
)

// GetRoot сервер должен отдавать HTML-страничку со списком имён и значений всех известных ему на текущий момент метрик.
func GetRoot(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	sd := StoredData
	if len(sd) != 0 {
		var value string
		for k := range sd {
			if sd[k].Mtype == "gauge" {
				value = strconv.FormatFloat(sd[k].Val.gauge, 'f', -1, 64)
			} else {
				value = strconv.FormatInt(sd[k].Val.counter, 10)
			}
			kvw := "type: " + sd[k].Mtype + " name: " + sd[k].Name + " value: " + value + "\n"
			_, err := w.Write([]byte(kvw))
			if err != nil {
				return
			}
		}
	} else {
		_, err := w.Write([]byte("Нет данных"))
		if err != nil {
			return
		}
	}
}

// PostMetrics читаем данные из URL и сохраняем
func PostMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	res := tools.GetURL(r.URL.Path, "update")

	data, answer := storeData(res)
	if !answer {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Что-то пошло не так"))
		if err != nil {
			return
		}
		return
	}
	StoredData = data
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Я всё записал"))
	if err != nil {
		return
	}
}

// GetValue должен возвращать текущее значение запрашиваемой метрики
// в текстовом виде по запросу GET
// http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
func GetValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	res := tools.GetURL(r.URL.Path, "value")
	metrics := []string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
		"RandomValue",
	}
	typeM := res[0]
	nameM := res[1]

	if !contains(metrics, nameM) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Нет такой метрики"))
		if err != nil {
			return
		}
		return
	} else if typeM != "gauge" && typeM != "counter" {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Нет такого типа метрики"))
		if err != nil {
			return
		}
		return
	} else {
		sd := StoredData
		var value string
		for k := range sd {
			if sd[k].Name == nameM && sd[k].Mtype == typeM {
				if typeM == "gauge" {
					value = strconv.FormatFloat(sd[k].Val.gauge, 'f', -1, 64)
				} else {
					value = strconv.FormatInt(sd[k].Val.counter, 10)
				}
				_, err := w.Write([]byte(value))
				if err != nil {
					return
				}
			}
		}
	}
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
