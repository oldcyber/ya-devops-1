package tools

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

//type configs interface {
//	GetDatabaseDSN() string
//}

func DBConnect(conn string) (*pgxpool.Pool, error) {
	var err error
	dbPool, err := pgxpool.Connect(context.Background(), conn)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer dbPool.Close()
	return dbPool, nil
}

// db ping
func Ping(ctx context.Context, dbpool *pgxpool.Pool) error {
	var err error
	_, err = dbpool.Exec(ctx, "SELECT 1")
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
