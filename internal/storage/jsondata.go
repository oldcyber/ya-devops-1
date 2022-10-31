//go:generate easyjson -all -snake_case $GOFILE

package storage

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/mailru/easyjson"
	"github.com/oldcyber/ya-devops-1/internal/env"
)

type Metrics struct {
	ID    string   `json:"id"`              // Имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // Значение хеш-функции
}

func (m *Metrics) MarshalGaugeMetrics(k, key string, v float64) []byte {
	m.ID = k
	m.MType = env.MetricGaugeType
	m.Value = &v
	if key != "" {
		// SHA256 hash
		h := hmac.New(sha256.New, []byte(key))
		// log.Infof("converting: %s:gauge:%f", k, v)
		fmt.Fprintf(h, "%s:gauge:%f", k, v)
		m.Hash = fmt.Sprintf("%x", h.Sum(nil))
		// log.Info("Hash gauge: ", m.Hash)
	}
	rawBytes, err := easyjson.Marshal(m)
	if err != nil {
		panic(err)
	}
	return rawBytes
}

func (m *Metrics) MarshalCounterMetrics(c int64, key string) []byte {
	m.ID = "PollCount"
	m.MType = env.MetricCounterType
	m.Delta = &c
	if key != "" {
		// SHA256 hash
		h := hmac.New(sha256.New, []byte(key))
		// log.Infof("converting:%s:counter:%d", m.ID, c)
		fmt.Fprintf(h, "%s:counter:%d", m.ID, c)
		m.Hash = fmt.Sprintf("%x", h.Sum(nil))
		// log.Info("Hash counter: ", m.Hash)
	}
	rawBytes, err := easyjson.Marshal(m)
	if err != nil {
		panic(err)
	}
	return rawBytes
}

func (m *Metrics) SendBulkMetrics(myMap map[string]float64) []byte {
	var rawBytes []byte

	metrics := make([]Metrics, 0)
	for k, v := range myMap {
		name := k
		val := v
		metrics = append(metrics, Metrics{
			ID:    name,
			MType: env.MetricGaugeType,
			Value: &val,
		})
	}

	rawBytes = append(rawBytes, '[')
	c := len(metrics)
	for i := range metrics {
		rawB, err := easyjson.Marshal(metrics[i])
		if err != nil {
			panic(err)
		}
		rawBytes = append(rawBytes, rawB...)
		if i < c-1 {
			rawBytes = append(rawBytes, ',')
		}
	}
	rawBytes = append(rawBytes, ']')

	return rawBytes
}
