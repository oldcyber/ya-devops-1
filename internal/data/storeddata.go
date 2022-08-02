package data

import (
	"net/http"
	"strconv"

	"github.com/mailru/easyjson"
	log "github.com/sirupsen/logrus"
)

type StoredDataIface interface {
	AddStoredData(res []string) (bool, int)
	GetStoredData() *map[string]string
	GetStoredDataByName(mtype, mname string)
	StoredDataToJSON() []Metrics
}

type StoredType struct {
	gauge   float64
	counter int64
	stype   string
}

type storedData struct {
	data map[string]StoredType
}

func NewstoredData() *storedData {
	return &storedData{}
}

func (s *storedData) StoreJSONToData(m Metrics) (int, []byte, error) {
	var (
		out    Metrics
		err    error
		result []byte
	)
	if s.data == nil {
		s.data = map[string]StoredType{}
	}
	switch m.MType {
	case "gauge":
		s.data[m.ID] = StoredType{gauge: *m.Value, stype: m.MType}
		out = Metrics{MType: "gauge", ID: m.ID, Value: m.Value, Delta: nil}
		result, err = easyjson.Marshal(out)
		if err != nil {
			return http.StatusBadRequest, nil, err
		}
		log.Println("Записали данные в метрику", m.ID, "значение", s.data[m.ID].gauge)
		return http.StatusOK, result, nil
	case "counter":
		tt := s.data[m.ID].counter
		*m.Delta += tt
		s.data[m.ID] = StoredType{counter: *m.Delta, stype: m.MType}
		out = Metrics{MType: "counter", ID: m.ID, Value: nil, Delta: m.Delta}
		result, err = easyjson.Marshal(out)
		if err != nil {
			return http.StatusBadRequest, nil, err
		}
		log.Println("Записали данные в метрику", m.ID, "значение", s.data[m.ID].counter)
		return http.StatusOK, result, nil
	default:
		return http.StatusBadRequest, nil, err
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
		tt.stype = res[0]
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
		s.data[res[1]] = StoredType{counter: t + c, stype: res[0]}
		log.Println("Записали данные в метрику", res[1], "значение", t+c)
		return true, 200
	default:
		return false, 400
	}
}

func (s storedData) GetStoredDataByParamToJSON(mtype, mname string) ([]byte, int) {
	var out Metrics
	var out1 Metrics
	var result []byte
	for i := range s.data {
		if i == mname {
			log.Println("Нашли данные по имени", mname, "которые совпадают с записью", i)
			if mtype == "gauge" {
				// if s.data[i].gauge != 0 {
				log.Println("Нашли данные тип", mtype, "значение", s.data[i].gauge)
				te := s.data[i].gauge
				out = Metrics{MType: "gauge", ID: i, Value: &te, Delta: nil}
				log.Println("Преобразовали данные в метрику", out)
				result, err := easyjson.Marshal(out)
				if err != nil {
					return nil, 400
				}
				return result, 200
				//}
			} else if mtype == "counter" {
				// if s.data[i].counter != 0 {
				log.Println("Нашли данные тип", mtype, "значение", s.data[i].counter)
				ce := s.data[i].counter
				out1 = Metrics{MType: "counter", ID: i, Value: nil, Delta: &ce}
				log.Println("Преобразовали данные в метрику", out)
				result, err := easyjson.Marshal(out1)
				if err != nil {
					log.Println(err)
					return nil, 400
				}
				return result, 200
				//}
			}
		}
	}
	log.Println("Не нашли данные по имени", mname)
	return result, 404
}

func (s storedData) StoredDataToJSON() []Metrics {
	var out Metrics
	var w []Metrics
	for k, v := range s.data {
		if v.stype == "gauge" {
			te := s.data[k].gauge
			out = Metrics{MType: "gauge", ID: k, Value: &te}
			w = append(w, out)
		} else if v.stype == "counter" {
			te := s.data[k].counter
			out = Metrics{MType: "counter", ID: k, Delta: &te}
			w = append(w, out)
		}
	}
	return w
}

func (s storedData) GetStoredData() map[string]string {
	r := make(map[string]string)
	for k, v := range s.data {
		if v.stype == "gauge" {
			r[k] = strconv.FormatFloat(v.gauge, 'f', -1, 64)
		} else if v.stype == "counter" {
			r[k] = strconv.FormatInt(v.counter, 10)
		}
	}
	return r
}

func (s storedData) GetStoredDataByName(mtype, mname string) (string, int) {
	log.Info("s.data", s.data)
	for i := range s.data {
		if i == mname {
			if mtype == "gauge" {
				return strconv.FormatFloat(s.data[i].gauge, 'f', -1, 64), http.StatusOK
			} else if mtype == "counter" {
				return strconv.FormatInt(s.data[i].counter, 10), http.StatusOK
			}
		}
	}
	return "", http.StatusNotFound
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
