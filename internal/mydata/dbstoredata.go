package mydata

import (
	"database/sql"

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

func (ms *dbStoreData) UpdateStoreDataItem(db *sql.DB, m Metrics) error {
	_, err := db.Exec("UPDATE metrics SET metric_gauge = $1, metric_counter = $2 WHERE metric_name = $3", m.Value, m.Delta, m.ID)
	if err != nil {
		log.Error(err)
		return err
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
	err := db.QueryRow("SELECT metric_name, metric_type, metric_gauge, metric_counter FROM metrics WHERE metric_name = $1 AND metric_type = $2", metricName, metricType).Scan(&storeData.MetricName, &storeData.MetricType, &storeData.MetricGauge, &storeData.MetricCounter)
	if err != nil {
		log.Error(err)
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

func FindStoreDataItem(db *sql.DB, metricName string) bool {
	var storeData dbStoreData
	err := db.QueryRow("SELECT * FROM metrics WHERE metric_name = $1", metricName).Scan(&storeData.MetricName, &storeData.MetricType, &storeData.MetricGauge, &storeData.MetricCounter)
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}
