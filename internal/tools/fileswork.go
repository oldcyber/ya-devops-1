package tools

import (
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type OutFile struct {
	file *os.File     // Файл для записи
	mtx  sync.RWMutex // Мьютекс для записи
}

const defaultFileMode = 0o755

func OpenWriteToFile(filename string, storeinterval time.Duration) (*OutFile, error) {
	var (
		file OutFile
		err  error
	)
	file.mtx.Lock()
	defer file.mtx.Unlock()

	if storeinterval == 0 {
		file.file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC|os.O_SYNC, defaultFileMode)
	} else {
		file.file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, defaultFileMode)
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info("Log file opened")
	return &OutFile{file: file.file}, nil
}

func (of *OutFile) WriteToFile(data []byte) error {
	of.mtx.Lock()
	defer of.mtx.Unlock()
	// Check if file is open
	if of.file == nil {
		return nil
	}
	_, err := of.file.Write(data)
	if err != nil {
		log.Error("Error writing to file: ", err)
		return err
	}
	log.Info("Записали в файл:", string(data))
	return nil
}

func (of *OutFile) CloseFile() error {
	// закрываем файл
	return of.file.Close()
}
