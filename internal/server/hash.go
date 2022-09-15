package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/oldcyber/ya-devops-1/internal/mydata"
	log "github.com/sirupsen/logrus"
)

func calcHash(cfg config, m mydata.Metrics) string {
	var d string
	// SHA256 hash
	h := hmac.New(sha256.New, []byte(cfg.GetKey()))
	switch m.MType {
	case "gauge":
		d = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	case "counter":
		d = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	}
	h.Write([]byte(d))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// checkHash Check incoming hash signature and compare it with stored hash
func checkHash(cfg config, m mydata.Metrics) bool {
	hash := calcHash(cfg, m)
	// log.Info("Input hash: ", m.Hash, " new hash: ", hash)
	if !hmac.Equal([]byte(m.Hash), []byte(hash)) {
		log.Info("Hash is not equal")
		return false
	}
	log.Info("Hash is equal")
	return true
}
