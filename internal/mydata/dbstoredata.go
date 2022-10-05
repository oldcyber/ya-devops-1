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

type DbStoreData struct {
	MetricName    string          `json:"metric_name"`
	MetricType    string          `json:"metric_type"`
	MetricGauge   sql.NullFloat64 `json:"metric_gauge,omitempty"`
	MetricCounter sql.NullInt64   `json:"metric_counter,omitempty"`
}

//type myDb struct {
//	db *sql.DB
//}

//func NewDBData() *DbStoreData {
//	return &DbStoreData{}
//}

//--------------------------------------------------------------
// Work with DB
//--------------------------------------------------------------

// Add a new metric to the store
func CreateStoreDataItem(db *sql.DB, m Metrics) error {
	_, err := db.Exec("INSERT INTO metrics (metric_name, metric_type, metric_gauge, metric_counter) VALUES ($1, $2, $3, $4)", m.ID, m.MType, m.Value, m.Delta)
	if err != nil {
		log.Error("Ошибка выполнения запроса на добавление: ", err)
		return err
	}
	return nil
}

// FindItemByMetricName - поиск метрики в хранилище
func FindItemByMetricName(db *sql.DB, metricName string) (DbStoreData, bool) {
	var storeData DbStoreData
	err := db.QueryRow("SELECT * FROM metrics WHERE metric_name = $1", metricName).Scan(&storeData.MetricName, &storeData.MetricType, &storeData.MetricGauge, &storeData.MetricCounter)
	if err != nil {
		log.Error(err)
		return storeData, false
	}
	return storeData, true
}

// FindItemByMetricNameANDMetricType - получение метрики из хранилища
func FindItemByMetricNameANDMetricType(db *sql.DB, metricName, metricType string) (DbStoreData, error) {
	var storeData DbStoreData
	err := db.QueryRow("SELECT metric_name, metric_type, metric_gauge, metric_counter FROM metrics WHERE metric_name = $1 AND metric_type = $2", metricName, metricType).Scan(&storeData.MetricName, &storeData.MetricType, &storeData.MetricGauge, &storeData.MetricCounter)
	if err != nil {
		log.Error("SELECT error: ", err)
		return storeData, err
	}
	return storeData, nil
}

// DeleteStoreDataItem - удаление метрики из хранилища
func (ms *DbStoreData) DeleteStoreDataItem(db *sql.DB, metricName string) error {
	_, err := db.Exec("DELETE FROM metrics WHERE metric_name = $1", metricName)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//--------------------------------------------------------------
// For Interface
//--------------------------------------------------------------

// UpdateStoreDataItem Обновление данных в БД (тип метрики. имя метрики, значение)
func UpdateStoreDataItem(db *sql.DB, mName, mType, mValue string) (DbStoreData, error) {
	var res DbStoreData
	switch mType {
	case "gauge":
		g, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			log.Error(err)
			return res, err
		}
		_, err = db.Exec("UPDATE metrics SET metric_gauge = $1 WHERE metric_name = $2", g, mName)
		if err != nil {
			log.Error(err)
			return res, err
		}
		err = db.QueryRow("SELECT metric_name, metric_type, metric_gauge, metric_counter FROM metrics WHERE metric_name = $1", mName).Scan(&res.MetricName, &res.MetricType, &res.MetricGauge, &res.MetricCounter)
		if err != nil {
			log.Error(err)
			return res, err
		}

	case "counter":
		c, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			log.Error(err)
			return res, err
		}
		log.Info("c: ", c)

		// Ищем старое значение
		data, result := FindItemByMetricName(db, mName)
		switch result {
		case true:
			c = c + data.MetricCounter.Int64
			sqlStatement := "UPDATE metrics SET metric_counter = $1 WHERE metric_name = $2;"
			_, err := db.Exec(sqlStatement, c, mName)
			if err != nil {
				log.Error(err)
				return res, err
			}
			err = db.QueryRow("SELECT metric_name, metric_type, metric_gauge, metric_counter FROM metrics WHERE metric_name = $1", mName).Scan(&res.MetricName, &res.MetricType, &res.MetricGauge, &res.MetricCounter)
			if err != nil {
				log.Error(err)
				return res, err
			}
		case false:
			log.Info("No data in DB. Create new record")
			err = CreateStoreDataItem(db, Metrics{ID: mName, MType: mType, Delta: &c})
			//_, err = db.Exec("INSERT INTO metrics (metric_name, metric_type, metric_counter) VALUES ($1, $2, $3)", mName, mType, c)
			if err != nil {
				log.Error(err)
				return res, err
			}
			err = db.QueryRow("SELECT metric_name, metric_type, metric_gauge, metric_counter FROM metrics WHERE metric_name = $1", mName).Scan(&res.MetricName, &res.MetricType, &res.MetricGauge, &res.MetricCounter)
			if err != nil {
				log.Error(err)
				return res, err
			}
		}
	}
	return res, nil
}

func (ms *DbStoreData) AddNewItemToDB(db *sql.DB, res []string) (bool, int) {
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
		// Записываем в БД
		_, err := UpdateStoreDataItem(db, res[0], res[1], res[2])
		if err != nil {
			log.Error(err)
			return false, 400
		}
		return true, 200
	case "counter":
		_, err := UpdateStoreDataItem(db, res[0], res[1], res[2])
		if err != nil {
			log.Error(err)
			return false, 400
		}
		return true, 200
	default:
		return false, 400
	}
}

func (ms *DbStoreData) AddNewItem(db *sql.DB, res []string) (bool, int) {
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
		// Записываем в БД
		_, err := UpdateStoreDataItem(db, res[0], res[1], res[2])
		if err != nil {
			log.Error(err)
			return false, 400
		}
		return true, 200
	case "counter":
		_, err := UpdateStoreDataItem(db, res[0], res[1], res[2])
		if err != nil {
			log.Error(err)
			return false, 400
		}
		return true, 200
	default:
		return false, 400
	}
}

func (ms *DbStoreData) GetStoredDataByNameFromDB(db *sql.DB, mType, mName string) (string, int) {
	// log.Info("ms.mydata", ms.MetricType)
	item, err := FindItemByMetricNameANDMetricType(db, mName, mType)
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

func (ms *DbStoreData) GetStoredDataByParamFromDBToJSON(db *sql.DB, m Metrics, key string) ([]byte, int) {
	var out Metrics
	var result []byte
	item, err := FindItemByMetricNameANDMetricType(db, m.ID, m.MType)
	if err != nil {
		log.Error("error find record: ", err)
		return nil, http.StatusNotFound
	}
	switch m.MType {
	case "gauge":
		log.Info("Нашли данные тип: ", m.MType, " значение: ", item.MetricGauge)
		te := item.MetricGauge.Float64
		hash := CountHash(key, "gauge", m.ID, te, 0)
		out = Metrics{MType: "gauge", ID: item.MetricName, Value: &te, Delta: nil, Hash: hash}
		log.Info("Преобразовали данные в метрику", out)
		result, err := easyjson.Marshal(out)
		if err != nil {
			return nil, http.StatusNotFound
		}
		return result, http.StatusOK
	case "counter":
		log.Info("Нашли данные тип: ", m.MType, " значение: ", item.MetricCounter)
		ce := item.MetricCounter.Int64
		hash := CountHash(key, "counter", m.ID, 0, ce)
		out = Metrics{MType: "counter", ID: item.MetricName, Value: nil, Delta: &ce, Hash: hash}
		log.Info("Преобразовали данные в метрику", out)
		result, err := easyjson.Marshal(out)
		if err != nil {
			log.Error(err)
			return nil, http.StatusNotFound
		}
		return result, http.StatusOK
	}

	log.Warn("Не нашли данные по имени", m.ID)
	return result, http.StatusNotFound
}

// StoreToDB store JSON to DB
func (ms DbStoreData) StoreToDB(db *sql.DB, m Metrics) (int, []byte, error) {
	var err error
	var out Metrics
	var result []byte

	_, res := FindItemByMetricName(db, m.ID)
	// log.Info("Search result: ", res)
	switch res {
	case false:
		// Запись не найдена - создаём новую запись
		err := CreateStoreDataItem(db, m)
		if err != nil {
			log.Error("Создание новой записи провалилось: ", err)
			return http.StatusBadRequest, nil, err
		}
		switch m.MType {
		case "gauge":
			out = Metrics{MType: m.MType, ID: m.ID, Value: m.Value, Delta: nil, Hash: m.Hash}
		case "counter":
			out = Metrics{MType: m.MType, ID: m.ID, Value: nil, Delta: m.Delta, Hash: m.Hash}
		}
		result, err = easyjson.Marshal(out)
		if err != nil {
			return http.StatusBadRequest, nil, err
		}
		return http.StatusOK, result, nil
	case true:
		// Запись найдена - обновляем значение
		switch m.MType {
		case "gauge":
			myRes, _ := UpdateStoreDataItem(db, m.ID, m.MType, fmt.Sprintf("%.12f", *m.Value))
			//if err != nil {
			//	log.Error("Обновление записи gauge провалилось: ", err)
			//	return http.StatusBadRequest, nil, err
			//}
			out = Metrics{MType: "gauge", ID: myRes.MetricName, Value: &myRes.MetricGauge.Float64, Delta: nil, Hash: m.Hash}
			result, err = easyjson.Marshal(out)
			if err != nil {
				return http.StatusBadRequest, nil, err
			}
			return http.StatusOK, result, nil
		case "counter":
			log.Info("*m.Delta: ", *m.Delta)
			log.Info("*m.Delta after convert: ", fmt.Sprintf("%d", *m.Delta))
			log.Info("*m.Delta after my convert: ", strconv.FormatInt(*m.Delta, 10))
			myRes, err := UpdateStoreDataItem(db, m.ID, m.MType, strconv.FormatInt(*m.Delta, 10))
			if err != nil {
				log.Error("Обновление записи counter провалилось: ", err)
				return http.StatusBadRequest, nil, err
			}
			out = Metrics{MType: "counter", ID: myRes.MetricName, Value: nil, Delta: &myRes.MetricCounter.Int64, Hash: m.Hash}
			result, err = easyjson.Marshal(out)
			if err != nil {
				return http.StatusBadRequest, nil, err
			}
			return http.StatusOK, result, nil
		default:
			return http.StatusBadRequest, nil, err
		}
	}
	return http.StatusBadRequest, nil, err
	// return http.StatusOK, nil, nil
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
