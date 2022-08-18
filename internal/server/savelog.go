package server

import (
	"time"

	"github.com/oldcyber/ya-devops-1/internal/mydata"
	"github.com/oldcyber/ya-devops-1/internal/tools"
	log "github.com/sirupsen/logrus"
)

type config interface {
	GetStoreFile() string
	GetStoreInterval() time.Duration
	GetKey() string
	CountHash(mydata.Metrics) string
	CheckHash(mydata.Metrics) bool
	GetDatabaseDSN() string
}

//type dbStoreData interface {
//	StoreJSONToDB(*sql.DB, mydata.Metrics) (int, []byte, error)
//}

//type outFile interface {
//	OpenWriteToFile(fileName string, interval time.Duration) (file *os.File, err error)
//}

func WorkWithLogs(cfg config) error {
	log.Info("Loading store file:", cfg.GetStoreFile(), " store interval:", cfg.GetStoreInterval())
	if cfg.GetStoreInterval() == 0 {
		log.Info("Надо писать в живую")
		return nil
	}

	timer1 := time.NewTicker(cfg.GetStoreInterval())
	defer func() {
		timer1.Stop()
	}()
	for {
		<-timer1.C
		log.Info("Start saving logs")
		f, err := tools.OpenWriteToFile(cfg.GetStoreFile(), cfg.GetStoreInterval())
		if err != nil {
			return err
		}
		err = SaveLog(f)
		if err != nil {
			return err
		}
		err = f.CloseFile()
		if err != nil {
			return err
		}
		log.Info("Log file saved")
	}
}
