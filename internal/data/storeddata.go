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

func (s *storedData) AddStoredJSONData(m *Metrics) (bool, int) {
	if s.data == nil {
		s.data = map[string]StoredType{}
	}
	log.Println("Начинаем запись данных", m.ID, m.MType, m.Value, m.Delta)
	switch m.MType {
	case "gauge":
		s.data[m.ID] = StoredType{gauge: *m.Value}
		return true, 200
	case "counter":
		tt := s.data[m.ID]
		log.Println("Предыдущее значение", tt.counter)
		tt.counter += *m.Delta
		log.Println("Новое значение", tt.counter)
		s.data[m.ID] = StoredType{counter: tt.counter}
		return true, 200
	default:
		return false, 400
	}
}

func (s *storedData) AddStoredData(res []string) (bool, int) {
	log.Println("Начинаем запись данных", len(res))
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

	log.Println("Проверяем тип метрики", res[0], "имя метрики", res[1], "значение", res[2])
	switch res[0] {
	case "gauge":
		g, err := strconv.ParseFloat(res[2], 64)
		if err != nil {
			log.Println(err)
			return false, 400
		}
		// Запись через присваивание
		tt := s.data[res[1]]
		tt.gauge = g
		s.data[res[1]] = tt
		log.Println("Записали данные в метрику", res[1], "значение", g)
		// s.data[res[1]] = StoredType{gauge: g}
		return true, 200
	case "counter":
		c, err := strconv.ParseInt(res[2], 10, 64)
		if err != nil {
			log.Println(err)
			return false, 400
		}
		tCounter := s.GetStoredData()
		t, _ := strconv.ParseInt(tCounter[res[1]], 10, 64)
		s.data[res[1]] = StoredType{counter: t + c}
		log.Println("Записали данные в метрику", res[1], "значение", t+c)
		return true, 200
	default:
		return false, 400
	}
}

func (s storedData) GetStoredDataByName(mtype, mname string) (string, int) {
	// log.Println("Начинаем поиск данных тип", mtype, "имя", mname)
	for i := range s.data {
		if i == mname {
			log.Println("Нашли данные по имени", mname, "которые совпадают с записью", i)
			if mtype == "gauge" {
				if s.data[i].gauge != 0 {
					log.Println("Нашли данные тип", mtype, "значение", s.data[i].gauge)
					return strconv.FormatFloat(s.data[i].gauge, 'f', -1, 64), 200
				}
			} else if mtype == "counter" {
				if s.data[i].counter != 0 {
					log.Println("Нашли данные тип", mtype, "значение", s.data[i].counter)
					return strconv.FormatInt(s.data[i].counter, 10), 200
				}
			}
		}
	}
	log.Println("Не нашли данные по имени", mname)
	return "", 404
}

func (s storedData) GetStoredData() map[string]string {
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
