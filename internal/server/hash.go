package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/oldcyber/ya-devops-1/internal/env"
	"github.com/oldcyber/ya-devops-1/internal/storage"
)

func CalcHash(cfg config, m storage.Metrics) string {
	var d string
	// SHA256 hash
	h := hmac.New(sha256.New, []byte(cfg.GetKey()))
	switch m.MType {
	case env.MetricGaugeType:
		d = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	case env.MetricCounterType:
		d = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	}
	h.Write([]byte(d))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// checkHash Check incoming hash signature and compare it with stored hash
func checkHash(cfg config, m storage.Metrics) bool {
	hash := CalcHash(cfg, m)
	return hmac.Equal([]byte(m.Hash), []byte(hash))
	//if !hmac.Equal([]byte(m.Hash), []byte(hash)) {
	//	// log.Info("Hash is not equal")
	//	return false
	//}
	//// log.Info("Hash is equal")
	//return true
}
