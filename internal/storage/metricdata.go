package storage

import (
	"crypto/rand"
	"math"
	"math/big"
	"runtime"
	"sync"

	log "github.com/sirupsen/logrus"
)

type (
	gauge   float64
	Counter int64
)

type MetricStore struct {
	data map[string]gauge
	mtx  sync.RWMutex
}

func NewMetricStore() *MetricStore {
	return &MetricStore{}
}

func (ms *MetricStore) GetMetrics() map[string]gauge {
	ms.mtx.RLock()
	defer ms.mtx.RUnlock()

	return ms.data
}

func (ms *MetricStore) AddMetrics(nm chan float64) {
	ms.mtx.RLock()
	defer ms.mtx.RUnlock()

	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	if ms.data == nil {
		ms.data = map[string]gauge{}
	}
	ms.data["Alloc"] = gauge(rtm.Alloc)
	ms.data["BuckHashSys"] = gauge(rtm.BuckHashSys)
	ms.data["Frees"] = gauge(rtm.Frees)
	ms.data["GCCPUFraction"] = gauge(rtm.GCCPUFraction)
	ms.data["GCSys"] = gauge(rtm.GCSys)
	ms.data["HeapAlloc"] = gauge(rtm.HeapAlloc)
	ms.data["HeapIdle"] = gauge(rtm.HeapIdle)
	ms.data["HeapInuse"] = gauge(rtm.HeapInuse)
	ms.data["HeapObjects"] = gauge(rtm.HeapObjects)
	ms.data["HeapReleased"] = gauge(rtm.HeapReleased)
	ms.data["HeapSys"] = gauge(rtm.HeapSys)
	ms.data["LastGC"] = gauge(rtm.LastGC)
	ms.data["Lookups"] = gauge(rtm.Lookups)
	ms.data["MCacheInuse"] = gauge(rtm.MCacheInuse)
	ms.data["MCacheSys"] = gauge(rtm.MCacheSys)
	ms.data["MSpanInuse"] = gauge(rtm.MSpanInuse)
	ms.data["MSpanSys"] = gauge(rtm.MSpanSys)
	ms.data["Mallocs"] = gauge(rtm.Mallocs)
	ms.data["NextGC"] = gauge(rtm.NextGC)
	ms.data["NumForcedGC"] = gauge(rtm.NumForcedGC)
	ms.data["NumGC"] = gauge(rtm.NumGC)
	ms.data["OtherSys"] = gauge(rtm.OtherSys)
	ms.data["PauseTotalNs"] = gauge(rtm.PauseTotalNs)
	ms.data["StackInuse"] = gauge(rtm.StackInuse)
	ms.data["StackSys"] = gauge(rtm.StackSys)
	ms.data["Sys"] = gauge(rtm.Sys)
	ms.data["TotalAlloc"] = gauge(rtm.TotalAlloc)
	ms.data["RandomValue"] = gauge(SetRandomValue())
	// Новые метрики
	ms.data["TotalMemory"] = gauge(<-nm)
	ms.data["FreeMemory"] = gauge(<-nm)
	ms.data["CPUutilization1"] = gauge(<-nm)
}

// SetRandomValue Генерируем случайное число
func SetRandomValue() float64 {
	res, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		log.Fatalln(err)
	}
	return float64(res.Int64())
}
