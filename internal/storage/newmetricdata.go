package storage

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	log "github.com/sirupsen/logrus"
)

func (ms *MetricStore) GetNewMetrics(ch chan<- float64) {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Error(err)
	}
	cu, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Error("cpu err: ", err)
	}
	tm := v.Total
	log.Info("Total Memory: ", tm)
	ch <- float64(tm)
	fm := v.Free
	log.Info("Free Memory: ", fm)
	ch <- float64(fm)
	if cu != nil {
		log.Info("CPU Util: ", cu)
		ch <- cu[0]
	} else {
		log.Info("CPU Util: ", 0)
		ch <- 0
	}
	close(ch)
}
