package storage

import (
	"encoding/json"
	"strconv"

	"github.com/mailru/easyjson"
	"github.com/oldcyber/ya-devops-1/internal/env"
	log "github.com/sirupsen/logrus"
)

type StoredData struct {
	StoredType
	s *StoredMem
}

type StoredType struct {
	gauge   float64
	counter int64
	stype   string
}

type StoredMem struct {
	data map[string]StoredType
}

func NewStoredMem() *StoredMem {
	return &StoredMem{
		data: map[string]StoredType{},
	}
}

func NewStoredData(sm *StoredMem) *StoredData {
	return &StoredData{
		s: sm,
	}
}

func (sd *StoredData) StoreTo(m Metrics) (code int, re []byte, er error) {
	var (
		out    Metrics
		err    error
		result []byte
	)
	if sd.s.data == nil {
		sd.s.data = map[string]StoredType{}
	}
	var tvalue float64
	var dvalue int64

	switch m.MType {
	case env.MetricGaugeType:
		if m.Value == nil {
			tvalue = 0.0
		} else {
			tvalue = *m.Value
		}
		sd.s.data[m.ID] = StoredType{gauge: tvalue, stype: m.MType}
		out = Metrics{ID: m.ID, MType: m.MType, Value: &tvalue, Hash: m.Hash}
	case env.MetricCounterType:
		if m.Delta == nil {
			dvalue = int64(0)
		} else {
			dvalue = *m.Delta
		}
		tt := sd.s.data[m.ID].counter
		// log.Info("Значение счетчика: ", tt)
		dvalue += tt
		sd.s.data[m.ID] = StoredType{counter: dvalue, stype: m.MType}
		out = Metrics{ID: m.ID, MType: m.MType, Delta: &dvalue, Hash: m.Hash}
	default:
		return StatusBadRequest, nil, err
	}
	result, err = json.Marshal(out)
	if err != nil {
		return StatusBadRequest, nil, err
	}
	return StatusOK, result, nil
}

func (sd *StoredData) AddNewItem(res []string) (status bool, code int) {
	// log.Info("Начинаем запись данных ", len(res))
	if sd.s.data == nil {
		sd.s.data = map[string]StoredType{}
	}
	// TODO: В функцию
	if len(res) < 3 {
		return false, StatusNotFound
	}
	types := []string{env.MetricGaugeType, env.MetricCounterType}

	if !contains(types, res[0]) {
		return false, StatusNotImplemented
	}

	// log.Info("Проверяем тип метрики: ", res[0], " имя метрики: ", res[1], " значение: ", res[2])
	switch res[0] {
	case env.MetricGaugeType:
		g, err := strconv.ParseFloat(res[2], env.BitSize)
		if err != nil {
			log.Error(err)
			return false, StatusBadRequest
		}
		// Запись через присваивание
		tt := sd.s.data[res[1]]
		tt.gauge = g
		tt.stype = res[0]
		sd.s.data[res[1]] = tt
		return true, StatusOK
	case env.MetricCounterType:
		c, err := strconv.ParseInt(res[2], env.Base, env.BitSize)
		if err != nil {
			log.Error(err)
			return false, StatusBadRequest
		}
		tCounter := sd.GetStoredData()
		t, _ := strconv.ParseInt(tCounter[res[1]], env.Base, env.BitSize)
		sd.s.data[res[1]] = StoredType{counter: t + c, stype: res[0]}
		return true, StatusOK
	default:
		return false, StatusBadRequest
	}
}

func (sd *StoredData) GetStoredDataByParamToJSON(m Metrics, key string) (re []byte, code int) {
	var out Metrics
	var result []byte
	var err error

	for i := range sd.s.data {
		if i == m.ID {
			// log.Info("Нашли данные по имени ", m.ID, " которые совпадают с записью ", i)
			if m.MType == env.MetricGaugeType {
				te := sd.s.data[i].gauge
				hash := CountHash(key, env.MetricGaugeType, m.ID, te, 0)
				out = Metrics{MType: env.MetricGaugeType, ID: i, Value: &te, Delta: nil, Hash: hash}
				result, err = easyjson.Marshal(out)
				if err != nil {
					return nil, StatusBadRequest
				}
				return result, StatusOK
			} else if m.MType == env.MetricCounterType {
				// log.Info("Нашли данные тип ", m.MType, " значение ", sd.s.data[i].counter)
				ce := sd.s.data[i].counter
				hash := CountHash(key, env.MetricCounterType, m.ID, 0, ce)
				out = Metrics{MType: env.MetricCounterType, ID: i, Value: nil, Delta: &ce, Hash: hash}
				result, err = easyjson.Marshal(out)
				if err != nil {
					log.Error(err)
					return nil, StatusBadRequest
				}
				return result, StatusOK
			}
		}
	}
	log.Warn("Не нашли данные по имени", m.ID)
	return result, StatusNotFound
}

func (sd *StoredData) GetDataToJSON() []Metrics {
	var out Metrics
	var w []Metrics
	for k, v := range sd.s.data {
		if v.stype == env.MetricGaugeType {
			te := sd.s.data[k].gauge
			out = Metrics{MType: env.MetricGaugeType, ID: k, Value: &te}
			w = append(w, out)
		} else if v.stype == env.MetricCounterType {
			te := sd.s.data[k].counter
			out = Metrics{MType: env.MetricCounterType, ID: k, Delta: &te}
			w = append(w, out)
		}
	}
	return w
}

func (sd *StoredData) GetStoredData() map[string]string {
	r := make(map[string]string)

	for k, v := range sd.s.data {
		if v.stype == env.MetricGaugeType {
			r[k] = strconv.FormatFloat(v.gauge, 'f', -prec, env.BitSize)
		} else if v.stype == env.MetricCounterType {
			r[k] = strconv.FormatInt(v.counter, env.Base)
		}
	}
	return r
}

func (sd *StoredData) GetStoredDataByName(mtype, mname string) (re string, code int) {
	for i := range sd.s.data {
		if i == mname {
			if mtype == env.MetricGaugeType {
				return strconv.FormatFloat(sd.s.data[i].gauge, 'f', -prec, env.BitSize), StatusOK
			} else if mtype == env.MetricCounterType {
				return strconv.FormatInt(sd.s.data[i].counter, env.Base), StatusOK
			}
		}
	}
	return "", StatusNotFound
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
