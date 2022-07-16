package server

import (
	"log"
	"net/http"
	"strconv"

	"ya-devops-1/internal/tools"
)

// GetRoot сервер должен отдавать HTML-страничку со списком имён и значений всех известных ему на текущий момент метрик.
func GetRoot(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	sd := StoredData
	for i := range sd {
		var ik string
		if sd[i].gauge != 0 {
			ik = "name: " + i + " value: " + strconv.FormatFloat(sd[i].gauge, 'f', -1, 64) + "\n"
		} else {
			ik = "name: " + i + " value: " + strconv.FormatInt(sd[i].counter, 10) + "\n"
		}
		_, err := w.Write([]byte(ik))
		if err != nil {
			return
		}
	}
}

// GetMetrics читаем данные из URL и сохраняем
func GetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	res := tools.GetURL(r.URL.Path, "update")
	if res == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
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

// GetValue должен возвращать текущее значение запрашиваемой метрики
// в текстовом виде по запросу GET
// http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
func GetValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	res := tools.GetURL(r.URL.Path, "value")
	typeM := res[0]
	nameM := res[1]

	if typeM != "gauge" && typeM != "counter" {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Нет такого типа метрики"))
		if err != nil {
			return
		}
		return
	} else {
		sd := StoredData
		for i := range sd {
			if i == nameM {
				if sd[i].gauge != 0 {
					value := strconv.FormatFloat(sd[i].gauge, 'f', -1, 64)
					_, err := w.Write([]byte(value))
					if err != nil {
						return
					}
				} else if sd[i].counter != 0 {
					value := strconv.FormatInt(sd[i].counter, 10)
					_, err := w.Write([]byte(value))
					if err != nil {
						return
					}
				} else {
					value := "HZ"
					_, err := w.Write([]byte(value))
					if err != nil {
						return
					}
				}
			}
		}
	}
}
