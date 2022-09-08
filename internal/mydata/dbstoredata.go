package mydata

import (
	"database/sql"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type dbStoreData struct {
	MetricName    string          `json:"metric_name"`
	MetricType    string          `json:"metric_type"`
	MetricGauge   sql.NullFloat64 `json:"metric_gauge,omitempty"`
	MetricCounter sql.NullInt64   `json:"metric_counter,omitempty"`
}

// Add a new metric to the store
func (ms *dbStoreData) CreateStoreDataItem(db *sql.DB, m Metrics) error {
	_, err := db.Exec("INSERT INTO metrics (metric_name, metric_type, metric_gauge, metric_counter) VALUES ($1, $2, $3, $4)", m.ID, m.MType, m.Value, m.Delta)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// UpdateStoreDataItem Обновление данных в БД (тиа метрики. имя метрики, значение)
func (ms *dbStoreData) UpdateStoreDataItem(db *sql.DB, mType, mName, mValue string) error {
	switch mType {
	case "gauge":
		g, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			log.Error(err)
		}
		_, err = db.Exec("UPDATE metrics SET metric_gauge = $1 WHERE metric_name = $2", g, mName)
		if err != nil {
			log.Error(err)
			return err
		}
	case "counter":
		c, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			log.Error(err)
		}
		// Ищем старое значение
		data, res := ms.FindStoreDataItem(db, mName)
		switch res {
		case true:
			c += data.MetricCounter.Int64
		}
		_, err = db.Exec("UPDATE metrics SET metric_counter = $1 WHERE metric_name = $2", c, mName)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func (ms *dbStoreData) DeleteStoreDataItem(db *sql.DB, metricName string) error {
	_, err := db.Exec("DELETE FROM metrics WHERE metric_name = $1", metricName)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (ms *dbStoreData) GetStoreDataItem(db *sql.DB, metricName, metricType string) (dbStoreData, error) {
	var storeData dbStoreData
	// log.Info("SELECT metric_name, metric_type, metric_gauge, metric_counter FROM metrics WHERE metric_name =" + metricName + " AND metric_type = " + metricType)
	err := db.QueryRow("SELECT metric_name, metric_type, metric_gauge, metric_counter FROM metrics WHERE metric_name = $1 AND metric_type = $2", metricName, metricType).Scan(&storeData.MetricName, &storeData.MetricType, &storeData.MetricGauge, &storeData.MetricCounter)
	if err != nil {
		log.Error("SELECT error: ", err)
		return storeData, err
	}
	return storeData, nil
}

//func GetAllStoreDataItems(db *sql.DB) ([]dbStoreData, error) {
//	var storeData []dbStoreData
//	rows, err := db.Query("SELECT metric_name, metric_type, metric_gauge, metric_counter FROM metrics")
//	if err != nil {
//		log.Error(err)
//		return storeData, err
//	}
//	defer rows.Close()
//	for rows.Next() {
//		var s dbStoreData
//		err := rows.Scan(&s.MetricName, &s.MetricType, &s.MetricGauge, &s.MetricCounter)
//		if err != nil {
//			log.Error(err)
//			return storeData, err
//		}
//		storeData = append(storeData, s)
//	}
//	return storeData, nil
//}

func (ms *dbStoreData) FindStoreDataItem(db *sql.DB, metricName string) (dbStoreData, bool) {
	var storeData dbStoreData
	err := db.QueryRow("SELECT * FROM metrics WHERE metric_name = $1", metricName).Scan(&storeData.MetricName, &storeData.MetricType, &storeData.MetricGauge, &storeData.MetricCounter)
	if err != nil {
		log.Error(err)
		return storeData, false
	}
	return storeData, true
}
