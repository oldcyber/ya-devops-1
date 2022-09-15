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

var (
	str   = mydata.NewstoredData() // cfg = tools.NewConfig()
	dbstr = mydata.NewDBData()
)

// ci  tools.Config
// ofi tools.OutFileInterface

type outFile interface {
	WriteToFile([]byte) error
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

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
		var m mydata.Metrics
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
		str.AddStoredData([]string{m.MType, m.ID, val})
	}
	return nil
}

func UpdateValue(h http.Handler, cfg config, db *sql.DB, myPing bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Работа с БД
		//var myPing bool
		//err := db.Ping()
		//if err != nil {
		//	myPing = false
		//	log.Error(err)
		//} else {
		//	myPing = true
		//}
		//defer db.Close()

		m := mydata.Metrics{}
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
		var status int
		var res []byte
		if !myPing {
			status, res, err = str.StoreToData(m)
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
		} else {
			status, res, err = dbstr.StoreToDB(db, m)
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

func GetValue(h http.Handler, cfg config, db *sql.DB, myPing bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Работа с БД
		//var myPing bool
		//err := db.Ping()
		//if err != nil {
		//	myPing = false
		//	log.Error(err)
		//} else {
		//	myPing = true
		//}
		//defer db.Close()

		var m mydata.Metrics
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
		var res []byte
		var status int
		if !myPing {
			res, status = str.GetStoredDataByParamToJSON(m, cfg.GetKey())
		} else {
			res, status = dbstr.GetStoredDBByParamToJSON(db, m, cfg.GetKey())
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

func GetDBMetric(h http.Handler, db *sql.DB, myPing bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		// Работа с БД
		//var myPing bool
		//err := db.Ping()
		//if err != nil {
		//	myPing = false
		//	log.Error(err)
		//} else {
		//	myPing = true
		//}
		//defer db.Close()

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

		var res string
		var status int
		if !myPing {
			res, status = str.GetStoredDataByName(typeM, nameM)
		} else {
			res, status = dbstr.GetStoredDBByName(db, typeM, nameM)
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

func UpdateDBMetrics(h http.Handler, cfg config, db *sql.DB, myPing bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		// Работа с БД
		//var myPing bool
		//err := db.Ping()
		//if err != nil {
		//	myPing = false
		//	log.Error(err)
		//} else {
		//	myPing = true
		//}
		//defer db.Close()

		var res []string
		res = append(res, chi.URLParam(r, "type"))
		res = append(res, chi.URLParam(r, "name"))
		res = append(res, chi.URLParam(r, "value"))

		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// dsn := cfg.GetDatabaseDSN()
		var er bool
		var an int
		if !myPing {
			er, an = str.AddStoredData(res)
			// er, an = dbstr.AddStoredDBData(db, res)
		} else {
			er, an = dbstr.AddStoredDBData(db, res)
			// er, an = str.AddStoredData(res)
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

func MassUpdateValues(h http.Handler, cfg config, db *sql.DB, myPing bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Работа с БД
		//var myPing bool
		//err := db.Ping()
		//if err != nil {
		//	myPing = false
		//	log.Error(err)
		//} else {
		//	myPing = true
		//}
		//defer db.Close()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("Ошибка в ReadAll", err)
			return
		}
		log.Info("BODY: ", string(body))
		var metrics []mydata.Metrics

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
			var status int
			var res []byte
			if !myPing {
				status, res, err = str.StoreToData(m)
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
			} else {
				status, res, err = dbstr.StoreToDB(db, m)
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
