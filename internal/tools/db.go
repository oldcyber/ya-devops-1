package tools

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

//type configs interface {
//	GetDatabaseDSN() string
//}

// var db *sql.DB

func DBConnect(conn string, ctx context.Context) (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	conn = conn + "&connect_timeout=10"
	log.Info("Connecting to database:", conn)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// defer dbPool.Close()
	return db, nil
}

// db ping
//func Ping(ctx context.Context) error {
//	var err error
//	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
//	defer cancel()
//
//	_, err = db.ExecContext(ctx, "SELECT pg_sleep(10)")
//	log.Error(err)
//	return err

//_, err = db.Exec("SELECT 1")
//if err != nil {
//	log.Error(err)
//	return err
//}
//return nil
//}
