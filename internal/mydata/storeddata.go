package mydata

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mailru/easyjson"
	log "github.com/sirupsen/logrus"
)

type StoredType struct {
	gauge   float64
	counter int64
	stype   string
}

type storedData struct {
	data map[string]StoredType
}

func NewstoredData() *storedData {
	return &storedData{}
}

func NewDBData() *dbStoreData {
	return &dbStoreData{}
}

// store JSON to DB
func (ms dbStoreData) StoreJSONToDB(db *sql.DB, m Metrics) (int, []byte, error) {
	_, res := ms.FindStoreDataItem(db, m.ID)
	log.Info("Search result: ", res)
	switch res {
	case false:
		// Запись не найдена - создаём новую запись
		err := ms.CreateStoreDataItem(db, m)
		if err != nil {
			log.Error(err)
			return 0, nil, err
		}
	case true:
		// Запись найдена - обновляем значение
		switch m.MType {
		case "gauge":
			err := ms.UpdateStoreDataItem(db, m.ID, m.MType, fmt.Sprintf("%f", *m.Value))
			if err != nil {
				log.Error(err)
				return 0, nil, err
			}
		case "counter":
			err := ms.UpdateStoreDataItem(db, m.ID, m.MType, strconv.FormatInt(*m.Delta, 10))
			if err != nil {
				log.Error(err)
				return 0, nil, err
			}
		}
	}
	return http.StatusOK, nil, nil
}

func (ms dbStoreData) GetStoredDBByParamToJSON(db *sql.DB, m Metrics, key string) ([]byte, int) {
	var out Metrics
	var result []byte
	item, err := ms.GetStoreDataItem(db, m.ID, m.MType)
	if err != nil {
		log.Error("error find record: ", err)
		return nil, 0
	}
	switch m.MType {
	case "gauge":
		log.Info("Нашли данные тип", m.MType, "значение", item.MetricGauge)
		te := item.MetricGauge.Float64
		hash := CountHash(key, "gauge", m.ID, te, 0)
		out = Metrics{MType: "gauge", ID: item.MetricName, Value: &te, Delta: nil, Hash: hash}
		log.Info("Преобразовали данные в метрику", out)
		result, err := easyjson.Marshal(out)
		if err != nil {
			return nil, 400
		}
		return result, 200
	case "counter":
		log.Info("Нашли данные тип", m.MType, "значение", item.MetricCounter)
		ce := item.MetricCounter.Int64
		hash := CountHash(key, "counter", m.ID, 0, ce)
		out = Metrics{MType: "counter", ID: item.MetricName, Value: nil, Delta: &ce, Hash: hash}
		log.Info("Преобразовали данные в метрику", out)
		result, err := easyjson.Marshal(out)
		if err != nil {
			log.Error(err)
			return nil, 400
		}
		return result, 200
	}

	log.Warn("Не нашли данные по имени", m.ID)
	return result, 404
}

func (ms dbStoreData) GetStoredDBByName(db *sql.DB, mType, mName string) (string, int) {
	// log.Info("ms.mydata", ms.MetricType)
	item, err := ms.GetStoreDataItem(db, mName, mType)
	if err != nil {
		return "", http.StatusNotFound
	}
	switch item.MetricType {
	case "gauge":
		return strconv.FormatFloat(item.MetricGauge.Float64, 'f', -1, 64), http.StatusOK
	case "counter":
		return strconv.FormatInt(item.MetricCounter.Int64, 10), http.StatusOK
	}

	return "", http.StatusNotFound
}

func (ms dbStoreData) AddStoredDBData(db *sql.DB, res []string) (bool, int) {
	//log.Info("Начинаем запись данных", len(res))
	//if s.data == nil {
	//	s.data = map[string]StoredType{}
	//}

	if len(res) < 3 {
		return false, 404
	}
	types := []string{"gauge", "counter"}

	if !contains(types, res[0]) {
		return false, 501
	}

	log.Info("Проверяем тип метрики: ", res[0], " имя метрики: ", res[1], " значение: ", res[2])
	switch res[0] {
	case "gauge":
		//g, err := strconv.ParseFloat(res[2], 64)
		//if err != nil {
		//	log.Error(err)
		//	return false, 400
		//}
		// Записываем в БД
		err := ms.UpdateStoreDataItem(db, res[0], res[1], res[2])
		if err != nil {
			log.Error(err)
			return false, 400
		}
		// Запись через присваивание
		// tt := ms.data[res[1]]
		// tt.gauge = g
		// tt.stype = res[0]
		// s.data[res[1]] = tt
		// log.Info("Записали данные в метрику", res[1], "значение", g)
		// s.mydata[res[1]] = StoredType{gauge: g}
		return true, 200
	case "counter":
		err := ms.UpdateStoreDataItem(db, res[0], res[1], res[2])
		if err != nil {
			log.Error(err)
			return false, 400
		}

		//c, err := strconv.ParseInt(res[2], 10, 64)
		//if err != nil {
		//	log.Error(err)
		//	return false, 400
		//}
		//tCounter := s.GetStoredData()
		//t, _ := strconv.ParseInt(tCounter[res[1]], 10, 64)
		//s.data[res[1]] = StoredType{counter: t + c, stype: res[0]}
		//log.Info("Записали данные в метрику", res[1], "значение", t+c)
		return true, 200
	default:
		return false, 400
	}
}

func (s *storedData) StoreJSONToData(m Metrics) (int, []byte, error) {
	var (
		out    Metrics
		err    error
		result []byte
	)
	if s.data == nil {
		s.data = map[string]StoredType{}
	}
	switch m.MType {
	case "gauge":
		s.data[m.ID] = StoredType{gauge: *m.Value, stype: m.MType}
		out = Metrics{MType: "gauge", ID: m.ID, Value: m.Value, Delta: nil, Hash: m.Hash}
		result, err = easyjson.Marshal(out)
		if err != nil {
			return http.StatusBadRequest, nil, err
		}
		// log.Info("Записали данные в метрику", m.ID, "значение", s.mydata[m.ID].gauge, "хэш", m.Hash)
		return http.StatusOK, result, nil
	case "counter":
		tt := s.data[m.ID].counter
		*m.Delta += tt
		s.data[m.ID] = StoredType{counter: *m.Delta, stype: m.MType}
		out = Metrics{MType: "counter", ID: m.ID, Value: nil, Delta: m.Delta, Hash: m.Hash}
		result, err = easyjson.Marshal(out)
		if err != nil {
			return http.StatusBadRequest, nil, err
		}
		// log.Info("Записали данные в метрику", m.ID, "значение", s.mydata[m.ID].counter, "хэш", m.Hash, "прирост", *m.Delta)
		return http.StatusOK, result, nil
	default:
		return http.StatusBadRequest, nil, err
	}
}

func (s *storedData) AddStoredData(res []string) (bool, int) {
	log.Info("Начинаем запись данных", len(res))
	if s.data == nil {
		s.data = map[string]StoredType{}
	}

	if len(res) < 3 {
		return false, 404
	}
	types := []string{"gauge", "counter"}

	if !contains(types, res[0]) {
		return false, 501
	}

	log.Info("Проверяем тип метрики", res[0], "имя метрики", res[1], "значение", res[2])
	switch res[0] {
	case "gauge":
		g, err := strconv.ParseFloat(res[2], 64)
		if err != nil {
			log.Error(err)
			return false, 400
		}
		// Запись через присваивание
		tt := s.data[res[1]]
		tt.gauge = g
		tt.stype = res[0]
		s.data[res[1]] = tt
		log.Info("Записали данные в метрику", res[1], "значение", g)
		// s.mydata[res[1]] = StoredType{gauge: g}
		return true, 200
	case "counter":
		c, err := strconv.ParseInt(res[2], 10, 64)
		if err != nil {
			log.Error(err)
			return false, 400
		}
		tCounter := s.GetStoredData()
		t, _ := strconv.ParseInt(tCounter[res[1]], 10, 64)
		s.data[res[1]] = StoredType{counter: t + c, stype: res[0]}
		log.Info("Записали данные в метрику", res[1], "значение", t+c)
		return true, 200
	default:
		return false, 400
	}
}

func (s storedData) GetStoredDataByParamToJSON(m Metrics, key string) ([]byte, int) {
	var out Metrics
	var result []byte

	for i := range s.data {
		if i == m.ID {
			log.Info("Нашли данные по имени ", m.ID, " которые совпадают с записью ", i)
			if m.MType == "gauge" {
				log.Info("Нашли данные тип ", m.MType, "значение ", s.data[i].gauge)
				te := s.data[i].gauge
				hash := CountHash(key, "gauge", m.ID, te, 0)
				out = Metrics{MType: "gauge", ID: i, Value: &te, Delta: nil, Hash: hash}
				log.Info("Преобразовали данные в метрику ", out)
				result, err := easyjson.Marshal(out)
				if err != nil {
					return nil, 400
				}
				return result, 200
			} else if m.MType == "counter" {
				log.Info("Нашли данные тип ", m.MType, " значение ", s.data[i].counter)
				ce := s.data[i].counter
				hash := CountHash(key, "counter", m.ID, 0, ce)
				out = Metrics{MType: "counter", ID: i, Value: nil, Delta: &ce, Hash: hash}
				log.Info("Преобразовали данные в метрику ", out)
				result, err := easyjson.Marshal(out)
				if err != nil {
					log.Error(err)
					return nil, 400
				}
				return result, 200
			}
			//switch m.MType {
			//case "gauge":
			//	log.Info("Нашли данные тип ", m.MType, "значение ", s.data[i].gauge)
			//	te := s.data[i].gauge
			//	hash := CountHash(key, "gauge", m.ID, te, 0)
			//	out = Metrics{MType: "gauge", ID: i, Value: &te, Delta: nil, Hash: hash}
			//	log.Info("Преобразовали данные в метрику ", out)
			//	result, err := easyjson.Marshal(out)
			//	if err != nil {
			//		return nil, 400
			//	}
			//	return result, 200
			//case "counter":
			//	log.Info("Нашли данные тип ", m.MType, " значение ", s.data[i].counter)
			//	ce := s.data[i].counter
			//	hash := CountHash(key, "counter", m.ID, 0, ce)
			//	out = Metrics{MType: "counter", ID: i, Value: nil, Delta: &ce, Hash: hash}
			//	log.Info("Преобразовали данные в метрику ", out)
			//	result, err := easyjson.Marshal(out)
			//	if err != nil {
			//		log.Error(err)
			//		return nil, 400
			//	}
			//	return result, 200
			//}
		}
	}
	log.Warn("Не нашли данные по имени", m.ID)
	return result, 404
}

func (s storedData) StoredDataToJSON() []Metrics {
	var out Metrics
	var w []Metrics
	for k, v := range s.data {
		if v.stype == "gauge" {
			te := s.data[k].gauge
			out = Metrics{MType: "gauge", ID: k, Value: &te}
			w = append(w, out)
		} else if v.stype == "counter" {
			te := s.data[k].counter
			out = Metrics{MType: "counter", ID: k, Delta: &te}
			w = append(w, out)
		}
	}
	return w
}

func (s storedData) GetStoredData() map[string]string {
	r := make(map[string]string)
	for k, v := range s.data {
		if v.stype == "gauge" {
			r[k] = strconv.FormatFloat(v.gauge, 'f', -1, 64)
		} else if v.stype == "counter" {
			r[k] = strconv.FormatInt(v.counter, 10)
		}
	}
	return r
}

func (s storedData) GetStoredDataByName(mtype, mname string) (string, int) {
	log.Info("s.mydata", s.data)
	for i := range s.data {
		if i == mname {
			if mtype == "gauge" {
				return strconv.FormatFloat(s.data[i].gauge, 'f', -1, 64), http.StatusOK
			} else if mtype == "counter" {
				return strconv.FormatInt(s.data[i].counter, 10), http.StatusOK
			}
		}
	}
	return "", http.StatusNotFound
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func CountHash(key, mtype, mid string, mvalue float64, mdelta int64) string {
	log.Info("Вычисляем хэш для метрики: ", mid, " типа: ", mtype, " mvalue: ", mvalue, " mdelta: ", mdelta)
	var d string
	// SHA256 hash
	h := hmac.New(sha256.New, []byte(key))
	switch mtype {
	case "gauge":
		d = fmt.Sprintf("%s:gauge:%f", mid, mvalue)
	case "counter":
		d = fmt.Sprintf("%s:counter:%d", mid, mdelta)
	}
	h.Write([]byte(d))
	return fmt.Sprintf("%x", h.Sum(nil))
}
