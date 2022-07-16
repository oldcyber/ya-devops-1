package server

import (
	"log"
	"strconv"
)

var StoredData map[int]*SData

type SData struct {
	Mtype string
	Name  string
	Val   storedType
}

type storedType struct {
	gauge   float64
	counter int64
}

//type stores interface {
//	storeData([]string) bool
//}

// storeData - хранит данные в памяти
func storeData(res []string) (map[int]*SData, bool) {
	var s *SData
	if len(res) != 3 {
		return StoredData, false
	}
	n, err := strconv.ParseFloat(res[2], 64)
	if err != nil {
		log.Println(err)
	}
	if res[0] == "gauge" {
		s = &SData{Mtype: res[0], Name: res[1], Val: storedType{gauge: n}}
	} else if res[0] == "counter" {
		c, err := strconv.Atoi(res[2])
		if err != nil {
			log.Println(err)
		}
		s = &SData{Mtype: res[0], Name: res[1], Val: storedType{counter: int64(c)}}
	} else {
		return StoredData, false
	}

	switch res[1] {
	case "Alloc":
		StoredData[0] = s
	case "BuckHashSys":
		StoredData[1] = s
	case "Frees":
		StoredData[2] = s
	case "GCCPUFraction":
		StoredData[3] = s
	case "GCSys":
		StoredData[4] = s
	case "HeapAlloc":
		StoredData[5] = s
	case "HeapIdle":
		StoredData[6] = s
	case "HeapInuse":
		StoredData[7] = s
	case "HeapObjects":
		StoredData[8] = s
	case "HeapReleased":
		StoredData[9] = s
	case "HeapSys":
		StoredData[10] = s
	case "LastGC":
		StoredData[11] = s
	case "Lookups":
		StoredData[12] = s
	case "MCacheInuse":
		StoredData[13] = s
	case "MCacheSys":
		StoredData[14] = s
	case "Mallocs":
		StoredData[15] = s
	case "NextGC":
		StoredData[16] = s
	case "NumForcedGC":
		StoredData[17] = s
	case "NumGC":
		StoredData[18] = s
	case "OtherSys":
		StoredData[19] = s
	case "PauseTotalNs":
		StoredData[20] = s
	case "StackInuse":
		StoredData[21] = s
	case "StackSys":
		StoredData[22] = s
	case "Sys":
		StoredData[23] = s
	case "TotalAlloc":
		StoredData[24] = s
	case "RandomValue":
		StoredData[25] = s
	case "PollCount":
		StoredData[26] = s
	}

	return StoredData, true
}
