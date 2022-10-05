package server

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/mailru/easyjson"

	log "github.com/sirupsen/logrus"

	"github.com/oldcyber/ya-devops-1/internal/data"

	"github.com/go-chi/chi/v5"
)

var str = data.NewstoredData() // cfg = tools.NewConfig()
// ci  tools.Config
// ofi tools.OutFileInterface

type outFile interface {
	WriteToFile([]byte) error
}

// GetRoot сервер должен отдавать HTML-страничку со списком имён и значений всех известных ему на текущий момент метрик.
func GetRoot(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-type")
	if contentType == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		res := str.StoredDataToJSON()
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
	} else {
		w.Header().Set("Content-Type", "text/html")
	}
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

// GzipMiddleware HTTP middleware setting a value on the request context
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}
		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			_, err := io.WriteString(w, err.Error())
			if err != nil {
				return
			}
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

// UpdateJSONMetrics читаем JSON из URL и сохраняем
func UpdateJSONMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	m := data.Metrics{}
	err := easyjson.UnmarshalFromReader(r.Body, &m)
	if err != nil {
		log.Println("Ошибка в Unmarshall", err)
		return
	}
	status, res, err := str.StoreJSONToData(m)
	if err != nil {
		w.WriteHeader(status)
		_, err = w.Write(res)
		if err != nil {
			log.Println("Ошибка в Write", err)
			return
		}
		log.Println(err)
		return
	} else {
		w.WriteHeader(status)
		_, err = w.Write(res)
		if err != nil {
			log.Println("Ошибка в Write", err)
			return
		}
	}
}

// UpdateMetrics читаем данные из URL и сохраняем
func UpdateMetrics(w http.ResponseWriter, r *http.Request) {
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

// GetMetric должен возвращать текущее значение запрашиваемой метрики
// в текстовом виде по запросу GET
// http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
func GetMetric(w http.ResponseWriter, r *http.Request) {
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
		log.Error("Ошибка в Write", err)
		return
	}
}

// GetJSONMetric должен возвращать текущее значение запрашиваемой метрики
// в текстовом виде по запросу GET
// http://<АДРЕС_СЕРВЕРА>/value/{JSON}
func GetJSONMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m data.Metrics
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		return
	}
	log.Info(m.ID, m.MType)
	typeM := m.MType
	nameM := m.ID
	if typeM != "gauge" && typeM != "counter" {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Нет такого типа метрики"))
		if err != nil {
			log.Error(err)
			return
		}
		return
	}

	res, status := str.GetStoredDataByParamToJSON(typeM, nameM)
	if status != 200 {
		w.WriteHeader(status)
		return
	}
	log.Println(string(res))
	_, err = w.Write(res)
	if err != nil {
		log.Error("Ошибка в Write", err)
		return
	}
}

// var OFile *tools.OutFile

// SaveLog is a function to save log to a file
func SaveLog(f outFile) error {
	// log.Info("Start function SaveLog")
	sdi := str.StoredDataToJSON()
	log.Info("sdi", sdi)
	for _, v := range sdi {
		marshal, err := easyjson.Marshal(v)
		// log.Info("marshal: ", string(marshal))
		marshal = append(marshal, '\n')
		if err != nil {
			return err
		}
		err = f.WriteToFile(marshal)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func ReadLogFile(cfg config) error {
	var val string
	log.Info("cfg.LogFile", cfg.GetStoreFile())
	fo, err := os.Open(cfg.GetStoreFile())
	if err != nil {
		log.Error(err)
		return err
	}
	defer fo.Close()

	scanner := bufio.NewScanner(fo)
	for scanner.Scan() {
		var m data.Metrics
		err := json.Unmarshal([]byte(scanner.Text()), &m)
		if err != nil {
			log.Error(err)
			return err
		}
		if m.MType == "gauge" {
			val = strconv.FormatFloat(*m.Value, 'f', -1, 64)
		} else if m.MType == "counter" {
			val = strconv.FormatInt(*m.Delta, 10)
		} else {
			log.Error("Нет такого типа метрики")
			return err
		}
		str.AddStoredData([]string{m.MType, m.ID, val})
	}
	return nil
}
