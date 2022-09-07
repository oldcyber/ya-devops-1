package tools

import (
	"database/sql"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func DBConnect(conn string) (*sql.DB, error) {
	// ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()
	conn += "&connect_timeout=10"
	log.Info("Connecting to database:", conn)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// defer dbPool.Close()
	return db, nil
}

func CreateTable(db *sql.DB) error {
	query := "CREATE TABLE IF NOT EXISTS metrics(metric_name varchar(255) NOT NULL, metric_type varchar(255) NOT NULL, metric_gauge double precision, metric_counter bigint, CONSTRAINT metrics_pkey PRIMARY KEY (metric_name))"
	// log.Info("Creating table:", query)
	_, err := db.Exec(query)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
