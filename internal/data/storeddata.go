package data

import (
	"log"
	"strconv"
)

type StoredDataIface interface {
	AddStoredData(res []string) (bool, int)
	GetStoredData() *map[string]string
	GetStoredDataByName(mtype, mname string)
	// GetStoredData()
}
type storedData struct {
	data map[string]StoredType
}

type StoredType struct {
	gauge   float64
	counter int64
}

func NewstoredData() *storedData {
	return &storedData{}
}

func (s *storedData) AddStoredData(res []string) (bool, int) {
	if s.data == nil {
		s.data = map[string]StoredType{}
	}

	if len(res) < 3 {
		return false, 404
	}
	types := []string{"gauge", "counter"}

	if !contains(types, res[0]) {
		return false, 501
	}

	switch res[0] {
	case "gauge":
		g, err := strconv.ParseFloat(res[2], 64)
		if err != nil {
			// log.Println(err)
			return false, 400
		}
		s.data[res[1]] = StoredType{gauge: g}
		return true, 200
	case "counter":
		c, err := strconv.ParseInt(res[2], 10, 64)
		if err != nil {
			// log.Println(err)
			return false, 400
		}
		tCounter := s.GetStoredData()
		t, _ := strconv.ParseInt(tCounter[res[1]], 10, 64)
		s.data[res[1]] = StoredType{counter: t + c}
		return true, 200
	default:
		return false, 400
	}
	//if res[0] == "gauge" {
	//	g, err := strconv.ParseFloat(res[2], 64)
	//	if err != nil {
	//		// log.Println(err)
	//		return false, 400
	//	}
	//	s.data[res[1]] = StoredType{gauge: g}
	//	return true, 200
	//}
	//if res[0] == "counter" {
	//	c, err := strconv.ParseInt(res[2], 10, 64)
	//	if err != nil {
	//		// log.Println(err)
	//		return false, 400
	//	}
	//	tCounter := s.GetStoredData()
	//	t, _ := strconv.ParseInt(tCounter[res[1]], 10, 64)
	//	s.data[res[1]] = StoredType{counter: t + c}
	//	return true, 200
	//}
	//return true, 200
}

func (s *storedData) GetStoredDataByName(mtype, mname string) (string, int) {
	log.Println("s.data", s.data)
	for i := range s.data {
		if i == mname {
			if mtype == "gauge" {
				return strconv.FormatFloat(s.data[i].gauge, 'f', -1, 64), 200
			} else if mtype == "counter" {
				return strconv.FormatInt(s.data[i].counter, 10), 200
			}
		}
	}
	return "", 404
}

func (s *storedData) GetStoredData() map[string]string {
	r := make(map[string]string)
	for k, v := range s.data {
		if v.gauge != 0 && v.counter == 0 {
			r[k] = strconv.FormatFloat(v.gauge, 'f', -1, 64)
		} else {
			r[k] = strconv.FormatInt(v.counter, 10)
		}
	}
	return r
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
