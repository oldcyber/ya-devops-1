package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/mailru/easyjson"

	log "github.com/sirupsen/logrus"

	"ya-devops-1/internal/data"

	"github.com/go-chi/chi/v5"
)

var str = data.NewstoredData()

// GetRoot сервер должен отдавать HTML-страничку со списком имён и значений всех известных ему на текущий момент метрик.
func GetRoot(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := str.GetStoredData()
	for _, v := range res {
		marshal, err := easyjson.Marshal(v)
		if err != nil {
			return
		}
		_, err = w.Write(marshal)
		if err != nil {
			return
		}
		_, err = w.Write([]byte("\n"))
		if err != nil {
			return
		}
	}
	//_, err := w.Write(res)
	//if err != nil {
	//	return
	//}
}

// GetJSONMetrics читаем JSON из URL и сохраняем
func GetJSONMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	m := data.Metrics{}
	b, err := io.ReadAll(r.Body)
	// err := easyjson.UnmarshalFromReader(r.Body, &m)
	if err != nil {
		log.Println("Ошибка в Body", err)
		return
	}
	log.Println("Полученные данные:", string(b))
	err = easyjson.Unmarshal(b, &m)
	if err != nil {
		log.Println("Ошибка в Unmarshall", err)
		return
	}
	// log.Println("Обработанные данные:", m)
	er, an := str.AddStoredJSONData(m)
	if !er {
		w.WriteHeader(an)
		return
	} else {
		w.WriteHeader(200)
	}
}

// GetMetrics читаем данные из URL и сохраняем
//func GetMetrics(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "text/plain")
//	var res []string
//	res = append(res, chi.URLParam(r, "type"))
//	res = append(res, chi.URLParam(r, "name"))
//	res = append(res, chi.URLParam(r, "value"))
//
//	if res == nil {
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//	er, an := str.AddStoredData(res)
//	if !er {
//		w.WriteHeader(an)
//		return
//	} else {
//		w.WriteHeader(200)
//	}
//}

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

// GetJSONValue должен возвращать текущее значение запрашиваемой метрики
// в текстовом виде по запросу GET
// http://<АДРЕС_СЕРВЕРА>/value/{JSON}
func GetJSONValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m data.Metrics
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		return
	}
	// ID и MType
	// log.Println(m)
	// log.Println(m.ID, m.MType)
	typeM := m.MType
	nameM := m.ID
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
	log.Println(string(res))
	_, err = w.Write(res)
	if err != nil {
		log.Println("Ошибка в Write", err)
		return
	}
}
