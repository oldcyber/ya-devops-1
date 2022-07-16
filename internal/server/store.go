package server

import (
	"strconv"
)

var (
	StoredData map[string]StoredType
	counter    = 0
)

type StoredType struct {
	gauge   float64
	counter int64
}

// storeData - хранит данные вида [string]gauge
func storeData(res []string) (bool, int) {
	if len(res) < 3 {
		return false, 404
	}
	types := []string{"gauge", "counter"}

	if !contains(types, res[0]) {
		return false, 501
	}

	if res[0] == "gauge" {
		g, err := strconv.ParseFloat(res[2], 64)
		if err != nil {
			return false, 501
		}
		StoredData[res[1]] = StoredType{gauge: g}
	} else if res[0] == "counter" {
		c, err := strconv.ParseInt(res[2], 10, 64)
		if err != nil {
			return false, 501
		}
		StoredData[res[1]] = StoredType{counter: c}
	}

	return true, 200
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
