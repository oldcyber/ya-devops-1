package data

import (
	"github.com/mailru/easyjson"
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
	rawBytes, err := easyjson.Marshal(m)
	// b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return rawBytes
}

func (m *Metrics) SendCounterMetrics(c int64) []byte {
	m.ID = "PollCount"
	m.MType = "counter"
	m.Delta = &c
	rawBytes, err := easyjson.Marshal(m)
	// b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return rawBytes
}
