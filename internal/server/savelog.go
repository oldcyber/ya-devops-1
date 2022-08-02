package server

import (
	"time"

	"ya-devops-1/internal/tools"

	log "github.com/sirupsen/logrus"
)

var (
	FO  *tools.OutFile
	err error
)

func WorkWithLogs() error {
	if tools.Conf.StoreInterval == 0 {
		log.Info("Надо писать в живую")
		return nil
	}
	// timeDuration := time.NewTicker(tools.Conf.StoreInterval)
	timer1 := time.NewTicker(tools.Conf.StoreInterval)
	defer func() {
		timer1.Stop()
	}()
	for {
		// select {
		// case
		<-timer1.C
		log.Info("Start saving logs")
		FO, err = tools.OpenWriteToFile(tools.Conf.StoreFile)
		if err != nil {
			log.Error(err)
			return err
		}
		err = SaveLog(FO)
		if err != nil {
			log.Error(err)
			return err
		}
		err = tools.CloseFile(FO)
		if err != nil {
			log.Error(err)
			return err
		}
		log.Info("Log file saved")
	}
	//}
}
