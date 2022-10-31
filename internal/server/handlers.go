package server

import (
	"bufio"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/mailru/easyjson"
	"github.com/oldcyber/ya-devops-1/internal/env"
	"github.com/oldcyber/ya-devops-1/internal/tools"
	log "github.com/sirupsen/logrus"

	"github.com/oldcyber/ya-devops-1/internal/storage"

	"github.com/go-chi/chi/v5"
)

type Storage interface {
	AddNewItem([]string) (bool, int)
	GetStoredDataByName(string, string) (string, int)
	StoreTo(storage.Metrics) (int, []byte, error)
	GetStoredDataByParamToJSON(storage.Metrics, string) ([]byte, int)
	GetDataToJSON() []storage.Metrics
}

type outFile interface {
	WriteToFile([]byte) error
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

const (
	saveToFile = "file"
	saveToDB   = "db"
)

// Ping при запросе проверяет соединение с базой данных.
// При успешной проверке хендлер должен вернуть HTTP-статус 200 OK, при неуспешной — 500 Internal Server Error.
func Ping(_ http.ResponseWriter, _ *http.Request) {
	// --------------------------------------------------------------
}

func GetPing(h http.Handler, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer db.Close()
		err := db.Ping()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		h.ServeHTTP(w, r)
	}
}

// GetRoot сервер должен отдавать HTML-страничку со списком имён и значений всех известных ему на текущий момент метрик.
func GetRoot(h http.Handler, db *sql.DB, ms *storage.StoredMem, storeTO string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-type")
		if contentType == "application/json" {
			w.Header().Set("Content-type", "application/json")
		} else {
			w.Header().Set("Content-type", "text/html")
		}

		var s Storage
		switch storeTO {
		case "file":
			s = storage.NewStoredData(ms)
		case saveToDB:
			s = storage.NewDBStoreData(db)
		}

		// log.Info("storeTO: ", storeTO)
		res := s.GetDataToJSON()
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
		h.ServeHTTP(w, r)
	}
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
			_, err = io.WriteString(w, err.Error())
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

// SaveLog is a function to save log to a file
func SaveLog(f outFile, ms *storage.StoredMem) error {
	var s Storage = storage.NewStoredData(ms)
	sdi := s.GetDataToJSON()
	// log.Info("sdi", sdi)
	for _, v := range sdi {
		marshal, err := easyjson.Marshal(v)
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

func ReadLogFile(cfg config, ms *storage.StoredMem) error {
	var val string
	// log.Info("cfg.LogFile: ", cfg.GetStoreFile())
	fo, err := os.Open(cfg.GetStoreFile())
	if err != nil {
		// Нет файла - создаем
		log.Info("Нет файла - создаем")
		_, err := tools.OpenWriteToFile(cfg.GetStoreFile(), cfg.GetStoreInterval())
		if err != nil {
			log.Error(err)
			return err
		}
	}
	defer fo.Close()

	scanner := bufio.NewScanner(fo)
	for scanner.Scan() {
		var (
			m *storage.Metrics
			s Storage
		)

		s = storage.NewStoredData(ms)
		err := json.Unmarshal([]byte(scanner.Text()), &m)
		if err != nil {
			log.Error(err)
			return err
		}
		switch {
		case m.MType == env.MetricGaugeType:
			val = strconv.FormatFloat(*m.Value, 'f', -1, env.BitSize)
		case m.MType == env.MetricCounterType:
			val = strconv.FormatInt(*m.Delta, env.Base)
		default:
			log.Error("Нет такого типа метрики")
			return err
		}
		s.AddNewItem([]string{m.MType, m.ID, val})
	}
	return nil
}

func StoreMetricsFromJSON(h http.Handler, cfg config, db *sql.DB, ms *storage.StoredMem, storeTO string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var (
			m      storage.Metrics
			status int
			res    []byte
			s      Storage
		)
		// log.Info("storeTO: ", storeTO)
		switch storeTO {
		case "file":
			s = storage.NewStoredData(ms)
		case saveToDB:
			s = storage.NewDBStoreData(db)
		}

		err := easyjson.UnmarshalFromReader(r.Body, &m)
		if err != nil {
			log.Error("Ошибка в Unmarshall: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if m.Hash != "" {
			if !checkHash(cfg, m) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		status, res, err = s.StoreTo(m)
		if err != nil {
			w.WriteHeader(status)
			_, err = w.Write(res)
			if err != nil {
				log.Error("Ошибка в Write", err)
				return
			}
			log.Error(err)
			return
		} else {
			w.WriteHeader(status)
			_, err = w.Write(res)
			if err != nil {
				log.Error("Ошибка в Write", err)
				return
			}
		}
		h.ServeHTTP(w, r)
	}
}

func GetMetricsFromJSON(h http.Handler, cfg config, db *sql.DB, ms *storage.StoredMem, storeTO string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var (
			m      storage.Metrics
			res    []byte
			status int
			s      Storage
		)

		switch storeTO {
		case saveToFile:
			s = storage.NewStoredData(ms)
		case saveToDB:
			s = storage.NewDBStoreData(db)
		}

		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			log.Error("Ошибка в Unmarshall", err)
			return
		}
		typeM := m.MType
		if typeM != env.MetricGaugeType && typeM != env.MetricCounterType {
			w.WriteHeader(http.StatusNotFound)
			_, err = w.Write([]byte("Нет такого типа метрики"))
			if err != nil {
				log.Error(err)
				return
			}
			return
		}

		res, status = s.GetStoredDataByParamToJSON(m, cfg.GetKey())

		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}
		// log.Info(string(res))
		_, err = w.Write(res)
		if err != nil {
			log.Error("Ошибка в Write", err)
			return
		}

		h.ServeHTTP(w, r)
	}
}

func GetMetrics(h http.Handler, db *sql.DB, ms *storage.StoredMem, storeTO string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		typeM := chi.URLParam(r, "type")
		nameM := chi.URLParam(r, "name")
		if typeM != env.MetricGaugeType && typeM != env.MetricCounterType {
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte("Нет такого типа метрики"))
			if err != nil {
				return
			}
			return
		}

		var (
			res    string
			status int
			s      Storage
		)

		switch storeTO {
		case saveToFile:
			s = storage.NewStoredData(ms)
		case saveToDB:
			s = storage.NewDBStoreData(db)
		}

		res, status = s.GetStoredDataByName(typeM, nameM)

		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}

		_, err := w.Write([]byte(res))
		if err != nil {
			log.Error("Ошибка в Write", err)
			return
		}
		h.ServeHTTP(w, r)
	}
}

func StoreMetrics(h http.Handler, db *sql.DB, ms *storage.StoredMem, storeTO string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		var (
			res []string
			er  bool
			an  int
			s   Storage
		)

		switch storeTO {
		case saveToFile:
			s = storage.NewStoredData(ms)
		case saveToDB:
			s = storage.NewDBStoreData(db)
		}
		res = append(res,
			chi.URLParam(r, "type"),
			chi.URLParam(r, "name"),
			chi.URLParam(r, "value"),
		)

		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		er, an = s.AddNewItem(res)

		if !er {
			w.WriteHeader(an)
			return
		} else {
			w.WriteHeader(http.StatusOK)
		}
		h.ServeHTTP(w, r)
	}
}

func MassStoreMetrics(h http.Handler, cfg config, db *sql.DB, ms *storage.StoredMem, storeTO string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("Ошибка в ReadAll", err)
			return
		}
		// log.Info("BODY: ", string(body))

		var (
			metrics []storage.Metrics
			status  int
			res     []byte
			s       Storage
		)

		switch storeTO {
		case saveToFile:
			s = storage.NewStoredData(ms)
		case saveToDB:
			s = storage.NewDBStoreData(db)
		}

		err = json.Unmarshal(body, &metrics)
		if err != nil {
			log.Error("Ошибка в Unmarshal 1 ", err)
			return
		}
		if metrics == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, m := range metrics {
			if m.Hash != "" {
				if !checkHash(cfg, m) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}

			status, res, err = s.StoreTo(m)
			if err != nil {
				w.WriteHeader(status)
				_, err = w.Write(res)
				if err != nil {
					log.Error("Ошибка в Write", err)
					return
				}
				log.Error(err)
				return
			} else {
				w.WriteHeader(status)
				_, err = w.Write(res)
				if err != nil {
					log.Error("Ошибка в Write", err)
					return
				}
			}
		}
		h.ServeHTTP(w, r)
	}
}

// Plug заглушка
func Plug(_ http.ResponseWriter, _ *http.Request) {
}
