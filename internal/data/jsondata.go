package data

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/mailru/easyjson"
	log "github.com/sirupsen/logrus"
)

type Metrics struct {
	ID    string   `json:"id"`              // Имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // Значение хеш-функции
}

func (m *Metrics) SendGaugeMetrics(k, key string, v float64) []byte {
	m.ID = k
	m.MType = "gauge"
	m.Value = &v
	if key != "" {
		// SHA256 hash
		h := hmac.New(sha256.New, []byte(key))
		log.Infof("converting: %s:gauge:%f", k, v)
		h.Write([]byte(fmt.Sprintf("%s:gauge:%f", k, v)))
		m.Hash = fmt.Sprintf("%x", h.Sum(nil))
		log.Info("Hash gauge: ", m.Hash)
	}
	rawBytes, err := easyjson.Marshal(m)
	if err != nil {
		panic(err)
	}
	return rawBytes
}

func (m *Metrics) SendCounterMetrics(c int64, key string) []byte {
	m.ID = "PollCount"
	m.MType = "counter"
	m.Delta = &c
	if key != "" {
		// SHA256 hash
		h := hmac.New(sha256.New, []byte(key))
		log.Infof("converting:%s:counter:%d", m.ID, c)
		h.Write([]byte(fmt.Sprintf("%s:counter:%d", m.ID, c)))
		m.Hash = fmt.Sprintf("%x", h.Sum(nil))
		log.Info("Hash counter: ", m.Hash)
	}
	rawBytes, err := easyjson.Marshal(m)
	if err != nil {
		panic(err)
	}
	return rawBytes
}
