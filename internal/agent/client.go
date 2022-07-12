package agent

import (
	"crypto/rand"
	"log"
	"math"
	"math/big"
	"runtime"
	"time"
)

// Объявляем переменные хранения
var (
	storedCounter counter = 0
	storedMetrics map[string]gauge
)

// Объявляем тип для метрик
type (
	gauge   float64
	counter int64
)

// WorkWithMetrics Начало клиента. Работа с таймерами
func WorkWithMetrics() {
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
			storedCounter++
			storedMetrics = getMetrics()
		case <-timer2.C:
			sendGaugeMetrics(storedMetrics)
			sendCounterMetrics(storedCounter)
		}
	}
}

// Собираем метрики
func getMetrics() map[string]gauge {
	m := make(map[string]gauge)
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	m["Alloc"] = gauge(rtm.Alloc)
	m["BuckHashSys"] = gauge(rtm.BuckHashSys)
	m["Frees"] = gauge(rtm.Frees)
	m["GCCPUFraction"] = gauge(rtm.GCCPUFraction)
	m["GCSys"] = gauge(rtm.GCSys)
	m["HeapAlloc"] = gauge(rtm.HeapAlloc)
	m["HeapIdle"] = gauge(rtm.HeapIdle)
	m["HeapInuse"] = gauge(rtm.HeapInuse)
	m["HeapObjects"] = gauge(rtm.HeapObjects)
	m["HeapReleased"] = gauge(rtm.HeapReleased)
	m["HeapSys"] = gauge(rtm.HeapSys)
	m["LastGC"] = gauge(rtm.LastGC)
	m["Lookups"] = gauge(rtm.Lookups)
	m["MCacheInuse"] = gauge(rtm.MCacheInuse)
	m["MCacheSys"] = gauge(rtm.MCacheSys)
	m["Mallocs"] = gauge(rtm.Mallocs)
	m["NextGC"] = gauge(rtm.NextGC)
	m["NumForcedGC"] = gauge(rtm.NumForcedGC)
	m["NumGC"] = gauge(rtm.NumGC)
	m["OtherSys"] = gauge(rtm.OtherSys)
	m["PauseTotalNs"] = gauge(rtm.PauseTotalNs)
	m["StackInuse"] = gauge(rtm.StackInuse)
	m["StackSys"] = gauge(rtm.StackSys)
	m["Sys"] = gauge(rtm.Sys)
	m["TotalAlloc"] = gauge(rtm.TotalAlloc)
	m["RandomValue"] = gauge(setRandomValue())
	return m
}

// Генерируем случайное число
func setRandomValue() float64 {
	res, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		log.Fatalln(err)
	}
	return float64(res.Int64())
}
