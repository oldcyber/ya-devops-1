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
	log "github.com/sirupsen/logrus"

	"github.com/oldcyber/ya-devops-1/internal/mydata"

	"github.com/go-chi/chi/v5"
)

type DBStorage interface {
	AddNewItemToDB(*sql.DB, []string) (bool, int)
	GetStoredDataByNameFromDB(*sql.DB, string, string) (string, int)
	StoreToDB(*sql.DB, mydata.Metrics) (int, []byte, error)
	GetStoredDataByParamFromDBToJSON(*sql.DB, mydata.Metrics, string) ([]byte, int)
}

type MapStorage interface {
	AddNewItemToFile([]string) (bool, int)
	GetStoredDataByName(string, string) (string, int)
	StoreToData(mydata.Metrics) (int, []byte, error)
	GetStoredDataByParamToJSON(mydata.Metrics, string) ([]byte, int)
	GetDataToJSON() []mydata.Metrics
}

type Storage interface {
	DBStorage
	MapStorage
}

type outFile interface {
	WriteToFile([]byte) error
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

//func NewHandlers(s Storage) *Handlers {
//	return &Handlers{s: s}
//}

// Ping при запросе проверяет соединение с базой данных.
// При успешной проверке хендлер должен вернуть HTTP-статус 200 OK, при неуспешной — 500 Internal Server Error.
func Ping(_ http.ResponseWriter, _ *http.Request) {
	// --------------------------------------------------------------
}

func GetPing(h http.Handler, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer db.Close()
		err := db.Ping()
		// err = tools.Ping(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		h.ServeHTTP(w, r)
	}
}

// GetRoot сервер должен отдавать HTML-страничку со списком имён и значений всех известных ему на текущий момент метрик.
func GetRoot(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-type")
	if contentType == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		var s Storage
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
	} else {
		w.Header().Set("Content-Type", "text/html")
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

// SaveLog is a function to save log to a file
func SaveLog(f outFile) error {
	var s Storage
	// log.Info("Start function SaveLog")
	sdi := s.GetDataToJSON()
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
		var (
			m mydata.Metrics
			s Storage
		)
		err := json.Unmarshal([]byte(scanner.Text()), &m)
		if err != nil {
			log.Error(err)
			return err
		}
		switch {
		case m.MType == "gauge":
			val = strconv.FormatFloat(*m.Value, 'f', -1, 64)
		case m.MType == "counter":
			val = strconv.FormatInt(*m.Delta, 10)
		default:
			log.Error("Нет такого типа метрики")
			return err
		}
		s.AddNewItemToFile([]string{m.MType, m.ID, val})
	}
	return nil
}

func StoreMetricsFromJSON(h http.Handler, cfg config, db *sql.DB, storeTO string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		m := mydata.Metrics{}

		var (
			status int
			res    []byte
			s      Storage
		)

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

		switch storeTO {
		case "db":
			status, res, err = s.StoreToDB(db, m)
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
		case "file":
			status, res, err = s.StoreToData(m)
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

func GetMetricsFromJSON(h http.Handler, cfg config, db *sql.DB, storeTO string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var (
			m      mydata.Metrics
			res    []byte
			status int
			s      Storage
		)
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			log.Error("Ошибка в Unmarshall", err)
			return
		}
		typeM := m.MType
		if typeM != "gauge" && typeM != "counter" {
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte("Нет такого типа метрики"))
			if err != nil {
				log.Error(err)
				return
			}
			return
		}

		switch storeTO {
		case "db":
			res, status = s.GetStoredDataByParamFromDBToJSON(db, m, cfg.GetKey())
		case "file":
			res, status = s.GetStoredDataByParamToJSON(m, cfg.GetKey())
		}

		if status != 200 {
			w.WriteHeader(status)
			return
		}
		log.Info(string(res))
		_, err = w.Write(res)
		if err != nil {
			log.Error("Ошибка в Write", err)
			return
		}

		h.ServeHTTP(w, r)
		// h(w, r)
	}
}

func GetMetrics(h http.Handler, db *sql.DB, storeTO string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var (
			res    string
			status int
			s      Storage
		)

		switch storeTO {
		case "db":
			res, status = s.GetStoredDataByNameFromDB(db, typeM, nameM)
		case "file":
			res, status = s.GetStoredDataByName(typeM, nameM)
		}

		if status != 200 {
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

func StoreMetrics(h http.Handler, db *sql.DB, storeTO string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		var (
			res []string
			er  bool
			an  int
			s   Storage
		)

		res = append(res, chi.URLParam(r, "type"))
		res = append(res, chi.URLParam(r, "name"))
		res = append(res, chi.URLParam(r, "value"))

		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		switch storeTO {
		case "db":
			er, an = s.AddNewItemToDB(db, res)
		case "file":
			er, an = s.AddNewItemToFile(res)
		}

		if !er {
			w.WriteHeader(an)
			return
		} else {
			w.WriteHeader(200)
		}
		h.ServeHTTP(w, r)
	}
}

func MassStoreMetrics(h http.Handler, cfg config, db *sql.DB, storeTO string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("Ошибка в ReadAll", err)
			return
		}
		log.Info("BODY: ", string(body))

		var (
			metrics []mydata.Metrics
			status  int
			res     []byte
			s       Storage
		)

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

			switch storeTO {
			case "db":
				status, res, err = s.StoreToDB(db, m)
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
			case "file":
				status, res, err = s.StoreToData(m)
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

		}
		h.ServeHTTP(w, r)
	}
}

// Plug заглушка
func Plug(_ http.ResponseWriter, _ *http.Request) {
}
