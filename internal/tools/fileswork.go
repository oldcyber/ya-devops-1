package tools

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type OutFile struct {
	file *os.File // Файл для записи
}

//func (o OutFile) Read(p []byte) (n int, err error) {
//	// TODO implement me
//
//	panic("implement me")
//}

func OpenWriteToFile(filename string) (*OutFile, error) {
	var (
		file *os.File
		err  error
	)
	if Conf.StoreInterval == 0 {
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC|os.O_SYNC, 0o755)
	} else {
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0o755)
	}

	if err != nil {
		log.Errorf("Error opening file: %v", err)
		return nil, err
	}
	log.Info("Log file opened")
	return &OutFile{file: file}, nil
}

func CloseFile(of *OutFile) error {
	// закрываем файл
	return of.file.Close()
}

func WriteToFile(of *OutFile, data []byte) error {
	_, err := of.file.Write(data)
	if err != nil {
		log.Errorf("Error writing to file: %v", err)
		return err
	}
	return nil
}

//func OpenReadFile(filename string) (*OutFile, error) {
//	file, err := os.OpenFile(filename, os.O_RDONLY, 0o755)
//	if err != nil {
//		log.Errorf("Error opening file: %v", err)
//		return nil, err
//	}
//	log.Info("Log file opened")
//	return &OutFile{file: file}, nil
//}
