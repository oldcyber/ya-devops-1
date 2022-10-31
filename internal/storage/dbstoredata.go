package storage

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/mailru/easyjson"
	"github.com/oldcyber/ya-devops-1/internal/env"
	log "github.com/sirupsen/logrus"
)

type DBStoreData struct {
	MetricName    string          `json:"metric_name"`
	MetricType    string          `json:"metric_type"`
	MetricGauge   sql.NullFloat64 `json:"metric_gauge,omitempty"`
	MetricCounter sql.NullInt64   `json:"metric_counter,omitempty"`
	db            *sql.DB
}

func (ms *DBStoreData) GetDataToJSON() []Metrics {
	// TODO implement me
	panic("implement me")
}

func NewDBStoreData(db *sql.DB) *DBStoreData {
	return &DBStoreData{
		db: db,
	}
}

const (
	StatusNotFound       = 404
	StatusNotImplemented = 501
	StatusBadRequest     = 400
	StatusOK             = 200
	prec                 = 1
)

//--------------------------------------------------------------
// Work with DB
//--------------------------------------------------------------

// Add a new metric to the store
func CreateStoreDataItem(db *sql.DB, m Metrics) error {
	_, err := db.Exec("INSERT INTO metrics (metric_name, metric_type, metric_gauge, metric_counter) VALUES"+
		" ($1, $2, $3, $4)", m.ID, m.MType, m.Value, m.Delta)
	if err != nil {
		log.Error("Ошибка выполнения запроса на добавление: ", err)
		return err
	}
	return nil
}

// FindItemByMetricName - поиск метрики в хранилище
func FindItemByMetricName(db *sql.DB, metricName string) (DBStoreData, bool) {
	var storeData DBStoreData
	err := db.QueryRow("SELECT * FROM metrics WHERE metric_name = $1", metricName).Scan(
		&storeData.MetricName,
		&storeData.MetricType,
		&storeData.MetricGauge,
		&storeData.MetricCounter)
	if err != nil {
		log.Error(err)
		return storeData, false
	}
	return storeData, true
}

// FindItemByMetricNameANDMetricType - получение метрики из хранилища
func FindItemByMetricNameANDMetricType(db *sql.DB, metricName, metricType string) (DBStoreData, error) {
	var storeData DBStoreData
	err := db.QueryRow("SELECT metric_name, metric_type, metric_gauge, metric_counter FROM"+
		" metrics WHERE metric_name = $1 AND metric_type = $2", metricName, metricType).Scan(
		&storeData.MetricName,
		&storeData.MetricType,
		&storeData.MetricGauge,
		&storeData.MetricCounter)
	if err != nil {
		log.Error("SELECT error: ", err)
		return storeData, err
	}
	return storeData, nil
}

// DeleteStoreDataItem - удаление метрики из хранилища
func (ms *DBStoreData) DeleteStoreDataItem(db *sql.DB, metricName string) error {
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
func UpdateStoreDataItem(db *sql.DB, mName, mType, mValue string) (DBStoreData, error) {
	var res DBStoreData
	switch mType {
	case env.MetricGaugeType:
		g, err := strconv.ParseFloat(mValue, env.BitSize)
		if err != nil {
			log.Error(err)
			return res, err
		}
		_, err = db.Exec("UPDATE metrics SET metric_gauge = $1 WHERE metric_name = $2", g, mName)
		if err != nil {
			log.Error(err)
			return res, err
		}
		err = db.QueryRow("SELECT metric_name, metric_type, metric_gauge, metric_counter FROM"+
			" metrics WHERE metric_name = $1", mName).Scan(
			&res.MetricName,
			&res.MetricType,
			&res.MetricGauge,
			&res.MetricCounter)
		if err != nil {
			log.Error(err)
			return res, err
		}

	case env.MetricCounterType:
		c, err := strconv.ParseInt(mValue, env.Base, env.BitSize)
		if err != nil {
			log.Error(err)
			return res, err
		}

		// Ищем старое значение
		data, result := FindItemByMetricName(db, mName)
		switch result {
		case true:
			c += data.MetricCounter.Int64
			sqlStatement := "UPDATE metrics SET metric_counter = $1 WHERE metric_name = $2;"
			_, err = db.Exec(sqlStatement, c, mName)
			if err != nil {
				log.Error(err)
				return res, err
			}
			err = db.QueryRow("SELECT metric_name, metric_type, metric_gauge, metric_counter FROM"+
				" metrics WHERE metric_name = $1", mName).Scan(
				&res.MetricName,
				&res.MetricType,
				&res.MetricGauge,
				&res.MetricCounter)
			if err != nil {
				log.Error(err)
				return res, err
			}
		case false:
			log.Info("No data in DB. Create new record")
			err = CreateStoreDataItem(db, Metrics{ID: mName, MType: mType, Delta: &c})
			if err != nil {
				log.Error(err)
				return res, err
			}
			err = db.QueryRow("SELECT metric_name, metric_type, metric_gauge, metric_counter FROM"+
				" metrics WHERE metric_name = $1", mName).Scan(
				&res.MetricName,
				&res.MetricType,
				&res.MetricGauge,
				&res.MetricCounter)
			if err != nil {
				log.Error(err)
				return res, err
			}
		}
	}
	return res, nil
}

func (ms *DBStoreData) AddNewItem(res []string) (status bool, code int) {
	if len(res) < 3 {
		return false, StatusNotFound
	}
	types := []string{env.MetricGaugeType, env.MetricCounterType}

	if !contains(types, res[0]) {
		return false, StatusNotImplemented
	}
	switch res[0] {
	case env.MetricGaugeType:
		// Записываем в БД
		_, err := UpdateStoreDataItem(ms.db, res[0], res[1], res[2])
		if err != nil {
			log.Error(err)
			return false, StatusBadRequest
		}
		return true, StatusOK
	case env.MetricCounterType:
		_, err := UpdateStoreDataItem(ms.db, res[0], res[1], res[2])
		if err != nil {
			log.Error(err)
			return false, StatusBadRequest
		}
		return true, StatusOK
	default:
		return false, StatusBadRequest
	}
}

func (ms *DBStoreData) GetStoredDataByName(mType, mName string) (res string, code int) {
	item, err := FindItemByMetricNameANDMetricType(ms.db, mName, mType)
	if err != nil {
		return "", StatusNotFound
	}
	switch item.MetricType {
	case env.MetricGaugeType:
		return strconv.FormatFloat(item.MetricGauge.Float64, 'f', -1, env.BitSize), StatusOK
	case env.MetricCounterType:
		return strconv.FormatInt(item.MetricCounter.Int64, env.Base), StatusOK
	}

	return "", StatusNotFound
}

func (ms *DBStoreData) GetStoredDataByParamToJSON(m Metrics, key string) (res []byte, code int) {
	var out Metrics
	var result []byte
	item, err := FindItemByMetricNameANDMetricType(ms.db, m.ID, m.MType)
	if err != nil {
		log.Error("error find record: ", err)
		return nil, StatusNotFound
	}
	switch m.MType {
	case env.MetricGaugeType:
		te := item.MetricGauge.Float64
		hash := CountHash(key, env.MetricGaugeType, m.ID, te, 0)
		out = Metrics{MType: env.MetricGaugeType, ID: item.MetricName, Value: &te, Delta: nil, Hash: hash}
		result, err = easyjson.Marshal(out)
		if err != nil {
			return nil, StatusNotFound
		}
		return result, StatusOK
	case env.MetricCounterType:
		ce := item.MetricCounter.Int64
		hash := CountHash(key, env.MetricCounterType, m.ID, 0, ce)
		out = Metrics{MType: env.MetricCounterType, ID: item.MetricName, Value: nil, Delta: &ce, Hash: hash}
		result, err = easyjson.Marshal(out)
		if err != nil {
			log.Error(err)
			return nil, StatusNotFound
		}
		return result, StatusOK
	}

	return result, StatusNotFound
}

// StoreTo store JSON to DB
func (ms DBStoreData) StoreTo(m Metrics) (code int, re []byte, er error) {
	var (
		err    error
		out    Metrics
		result []byte
		myRes  DBStoreData
	)
	_, res := FindItemByMetricName(ms.db, m.ID)
	switch res {
	case false:
		// Запись не найдена - создаём новую запись
		err = CreateStoreDataItem(ms.db, m)
		if err != nil {
			log.Error("Создание новой записи провалилось: ", err)
			return StatusBadRequest, nil, err
		}
		switch m.MType {
		case env.MetricGaugeType:
			out = Metrics{MType: m.MType, ID: m.ID, Value: m.Value, Delta: nil, Hash: m.Hash}
		case env.MetricCounterType:
			out = Metrics{MType: m.MType, ID: m.ID, Value: nil, Delta: m.Delta, Hash: m.Hash}
		}
		result, err = easyjson.Marshal(out)
		if err != nil {
			return StatusBadRequest, nil, err
		}
		return StatusOK, result, nil
	case true:
		// Запись найдена - обновляем значение
		switch m.MType {
		case env.MetricGaugeType:
			myRes, _ = UpdateStoreDataItem(ms.db, m.ID, m.MType, fmt.Sprintf("%.12f", *m.Value))
			out = Metrics{MType: env.MetricGaugeType, ID: myRes.MetricName, Value: &myRes.MetricGauge.Float64, Delta: nil, Hash: m.Hash}
			result, err = easyjson.Marshal(out)
			if err != nil {
				return StatusBadRequest, nil, err
			}
			return StatusOK, result, nil
		case env.MetricCounterType:
			myRes, err = UpdateStoreDataItem(ms.db, m.ID, m.MType, strconv.FormatInt(*m.Delta, env.Base))
			if err != nil {
				log.Error("Обновление записи counter провалилось: ", err)
				return StatusBadRequest, nil, err
			}
			out = Metrics{MType: env.MetricCounterType, ID: myRes.MetricName, Value: nil, Delta: &myRes.MetricCounter.Int64, Hash: m.Hash}
			result, err = easyjson.Marshal(out)
			if err != nil {
				return StatusBadRequest, nil, err
			}
			return StatusOK, result, nil
		default:
			return StatusBadRequest, nil, err
		}
	}
	return StatusBadRequest, nil, err
}

func CountHash(key, mtype, mid string, mvalue float64, mdelta int64) string {
	var d string
	// SHA256 hash
	h := hmac.New(sha256.New, []byte(key))
	switch mtype {
	case env.MetricGaugeType:
		d = fmt.Sprintf("%s:gauge:%f", mid, mvalue)
	case env.MetricCounterType:
		d = fmt.Sprintf("%s:counter:%d", mid, mdelta)
	}
	h.Write([]byte(d))
	return fmt.Sprintf("%x", h.Sum(nil))
}
