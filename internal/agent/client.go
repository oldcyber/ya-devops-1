package agent

import (
	"time"

	"ya-devops-1/internal/data"
)

// WorkWithMetrics Начало клиента. Работа с таймерами
func WorkWithMetrics() {
	c := data.NewCounter()
	m := data.NewMetricStore()
	timer1 := time.NewTicker(2 * time.Second)
	// mutex
	timer2 := time.NewTicker(10 * time.Second)
	defer func() {
		timer1.Stop()
		timer2.Stop()
	}()
	for {
		select {
		case <-timer1.C:
			c.IncCounter()
			m.AddMetrics()
		case <-timer2.C:
			r := make(map[string]float64)
			for key, val := range m.GetMetrics() {
				r[key] = float64(val)
			}
			sendGaugeMetrics(r)
			sendCounterMetrics(int64(c.Count()))
		}
	}
}
