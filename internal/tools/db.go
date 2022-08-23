package tools

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func DBConnect(conn string, ctx context.Context) (*sql.DB, error) {
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
