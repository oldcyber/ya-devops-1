package tools

import (
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type outFile struct {
	file *os.File     // Файл для записи
	mtx  sync.RWMutex // Мьютекс для записи
}

//func NewOutFile() *outFile {
//	return &outFile{
//		file: nil,
//		mtx:  sync.RWMutex{},
//	}
//}

func OpenWriteToFile(filename string, storeinterval time.Duration) (*outFile, error) {
	var (
		file outFile
		err  error
	)
	file.mtx.Lock()
	defer file.mtx.Unlock()

	if storeinterval == 0 {
		file.file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC|os.O_SYNC, 0o755)
	} else {
		file.file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0o755)
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info("Log file opened")
	return &outFile{file: file.file}, nil
}

func (of *outFile) WriteToFile(data []byte) error {
	of.mtx.Lock()
	defer of.mtx.Unlock()
	_, err := of.file.Write(data)
	if err != nil {
		log.Errorf("Error writing to file: %v", err)
		return err
	}
	return nil
}

func (of *outFile) CloseFile() error {
	// закрываем файл
	return of.file.Close()
}
