package server

import (
	"fmt"
	"strconv"
)

// storeData - хранит данные вида [string]gauge
func storeData(res []string) (bool, error) {
	if len(res) != 3 {
		return false, fmt.Errorf("wrong format of data")
	}
	storedData = make(map[string]gauge)
	n, err := strconv.ParseInt(res[2], 10, 64)
	if err == nil {
		fmt.Printf("%d of type %T", n, n)
	}
	storedData[res[1]] = gauge(n)
	return true, nil
}
