package tools

import (
	"database/sql"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func DBConnect(conn string) (*sql.DB, error) {
	conn += "&connect_timeout=10"
	log.Info("Connecting to database: ", conn)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Error("sql.Open error: ", err)
		return nil, err
	}
	return db, nil
}

func CreateTable(db *sql.DB) error {
	query := "CREATE TABLE IF NOT EXISTS metrics(metric_name varchar(255) NOT NULL, metric_type" +
		" varchar(255) NOT NULL, metric_gauge double precision, metric_counter bigint, CONSTRAINT metrics_pkey PRIMARY KEY (metric_name))"
	_, err := db.Exec(query)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
