package internal

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"time"
	"ya-devops-1/internal/model"
)

var counter = 0

func sendURL(req *http.Request) {
	client := &http.Client{}
	req.Header.Add("Content-Type", "text/plain")
	_, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		//os.Exit(1)
	}
	//fmt.Println("Статус-код ", response.Status)
	//defer response.Body.Close()
}
func SendResult(res *model.Metrics) {
	// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	endpoint := "http://127.0.0.1:8080/update/"
	//client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.Alloc).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.Alloc)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.BuckHashSys).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.BuckHashSys)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.Frees).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.Frees)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.GCCPUFraction).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.GCCPUFraction)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.GCSys).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.GCSys)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.HeapAlloc).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.HeapAlloc)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.HeapIdle).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.HeapIdle)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.HeapInuse).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.HeapInuse)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.HeapObjects).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.HeapObjects)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.HeapReleased).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.HeapReleased)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.HeapSys).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.HeapSys)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.LastGC).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.LastGC)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.Lookups).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.Lookups)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.MCacheInuse).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.MCacheInuse)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.MCacheSys).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.MCacheSys)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.MSpanInuse).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.MSpanInuse)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.MSpanSys).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.MSpanSys)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.Mallocs).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.Mallocs)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.NextGC).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.NextGC)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.NumForcedGC).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.NumForcedGC)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.NumGC).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.NumGC)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.OtherSys).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.OtherSys)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.PauseTotalNs).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.PauseTotalNs)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.StackInuse).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.StackInuse)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.StackSys).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.StackSys)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.Sys).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.Sys)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.TotalAlloc).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.TotalAlloc)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.PollCount).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.PollCount)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)
	req, err = http.NewRequest(http.MethodPost, endpoint+reflect.TypeOf(res.RandomValue).Name()+"/"+"Alloc"+"/"+fmt.Sprintf("%v", reflect.ValueOf(res.RandomValue)), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendURL(req)

}

func StoreMetrics() {
	var m = model.New()
	var res *model.Metrics
	timer1 := time.NewTicker(2 * time.Second)
	timer2 := time.NewTicker(10 * time.Second)
	defer func() {
		timer1.Stop()
		timer2.Stop()
	}()
	for {
		select {
		case <-timer1.C:
			m.PollCount++
			res = m.Get()
			counter++
			fmt.Println("StoreMetrics", res)
		case <-timer2.C:
			SendResult(res)
		}
	}

}
