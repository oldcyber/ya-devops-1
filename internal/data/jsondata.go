package data

import (
	"encoding/json"
	"log"
)

type Metrics struct {
	ID    string   `json:"id"`              // Имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
}

func (m *Metrics) SendGaugeMetrics(k string, v float64) []byte {
	m.ID = k
	m.MType = "gauge"
	m.Value = &v
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return b
}

func (m *Metrics) SendCounterMetrics(c int64) []byte {
	m.ID = "PollCount"
	m.MType = "counter"
	m.Delta = &c
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	log.Println(string(b))
	return b
}
