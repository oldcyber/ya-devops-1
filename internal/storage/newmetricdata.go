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
	c, _ := cpu.Counts(true)
	log.Info("CPU: ", c)
	ci, _ := cpu.Info()
	log.Info("CPU Info: ", ci)
	cu, _ := cpu.Percent(time.Second, true)
	tm := v.Total
	log.Info("Total Memory: ", tm)
	ch <- float64(tm)
	fm := v.Free
	log.Info("Free Memory: ", fm)
	ch <- float64(fm)
	log.Info("CPU Util: ", cu)
	if cu != nil {
		for _, i := range cu {
			ch <- i
		}
	} else {
		ch <- 0
	}

	close(ch)
}
