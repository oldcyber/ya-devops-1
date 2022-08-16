package server

import (
	"time"

	"github.com/oldcyber/ya-devops-1/internal/tools"
	log "github.com/sirupsen/logrus"
)

type config interface {
	GetStoreFile() string
	GetStoreInterval() time.Duration
}

//type outFile interface {
//	OpenWriteToFile(fileName string, interval time.Duration) (file *os.File, err error)
//}

func WorkWithLogs(cfg config) error {
	log.Info("Loading store file:", cfg.GetStoreFile(), " store interval:", cfg.GetStoreInterval())
	f, err := tools.OpenWriteToFile(cfg.GetStoreFile(), cfg.GetStoreInterval())
	//_, err := outFile.OpenWriteToFile(cfg.GetStoreFile(), cfg.GetStoreInterval())
	if err != nil {
		log.Error(err)
		return err
	}
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
		err = SaveLog(f)
		if err != nil {
			log.Error(err)
			return err
		}
		err = f.CloseFile()
		if err != nil {
			log.Error(err)
			return err
		}
		log.Info("Log file saved")
	}
	//}
}
