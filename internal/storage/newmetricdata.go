package storage

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	log "github.com/sirupsen/logrus"
)

func GetNewMetrics(ch chan []gauge) {
	var test []gauge
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Error(err)
	}

	cu, _ := cpu.Percent(0, false)
	//log.Info("Total CPU Util: ", len(cu))
	//for _, i := range cu {
	//	log.Info(i)
	//	// tm += uint64(i)
	//}
	tm := v.Total
	fm := v.Free
	log.Info("CPU: ", cu[0])
	test = append(test, gauge(tm), gauge(fm), 10, gauge(cu[0]))
	ch <- test
}
