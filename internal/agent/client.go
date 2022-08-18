package agent

import (
	"time"

	"github.com/oldcyber/ya-devops-1/internal/data"
)

type config interface {
	GetPollInterval() time.Duration
	GetReportInterval() time.Duration
	GetAddress() string
	GetRestore() bool
	GetKey() string
}

func WorkWithMetrics(cfg config) error {
	c := data.NewCounter()
	m := data.NewMetricStore()
	timer1 := time.NewTicker(cfg.GetPollInterval())
	timer2 := time.NewTicker(cfg.GetReportInterval())

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
			sendJSONGaugeMetrics(r, cfg)
			sendJSONCounterMetrics(int64(c.Count()), cfg)
		}
	}
}
