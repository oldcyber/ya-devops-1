package data

import (
	"reflect"
	"testing"
)

func Test_storedData_AddStoredJSONData(t *testing.T) {
	type fields struct {
		data map[string]StoredType
	}

	type args struct {
		m Metrics
	}
	var g float64 = 1.034
	var c int64 = 2
	tests := []struct {
		name   string
		fields fields
		args   args
		err    error
		status int
		res    []byte
	}{
		{
			name: "Проверка на правильность записи gauge данных в метрику",
			fields: fields{
				data: map[string]StoredType{
					"TestGauge":   {gauge: g},
					"TestCounter": {counter: 2},
				},
			},
			args: args{
				m: Metrics{
					ID:    "TestGauge",
					MType: "gauge",
					Value: &g,
				},
			},
			err:    nil,
			status: 200,
			res:    []byte(`{"id":"TestGauge","type":"gauge","value":1.034}`),
		},
		{
			name: "Проверка на правильность записи counter данных в метрику",
			fields: fields{
				data: map[string]StoredType{
					"TestGauge":   {gauge: g},
					"TestCounter": {counter: c},
				},
			},
			args: args{
				m: Metrics{
					ID:    "TestCounter",
					MType: "counter",
					Delta: &c,
				},
			},
			err:    nil,
			status: 200,
			res:    []byte(`{"id":"TestCounter","type":"counter","delta":4}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storedData{
				data: tt.fields.data,
			}
			err, status, res := s.StoreJSONToData(tt.args.m)
			if err != tt.err {
				t.Errorf("StoreJSONToData() err = %v, err %v", err, tt.err)
			}
			if status != tt.status {
				t.Errorf("StoreJSONToData() status = %v, err %v", status, tt.status)
			}
			if string(res) != string(tt.res) {
				t.Errorf("StoreJSONToData() res = %v, err %v", string(res), string(tt.res))
			}
		})
	}
}

func Test_storedData_GetStoredDataByName(t *testing.T) {
	type fields struct {
		data map[string]StoredType
	}
	type args struct {
		mtype string
		mname string
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		body       string
		statuscode int
	}{
		{
			name: "Проверка на правильность получения gauge данных из метрики",
			fields: fields{
				data: map[string]StoredType{
					"TestGauge":   {gauge: 1.034},
					"TestCounter": {counter: 2},
				},
			},
			args: args{
				mtype: "gauge",
				mname: "TestGauge",
			},
			body:       `{"id":"TestGauge","type":"gauge","value":1.034}`,
			statuscode: 200,
		},
		{
			name: "Проверка на правильность получения counter данных из метрики",
			fields: fields{
				data: map[string]StoredType{
					"TestGauge":   {gauge: 1.034},
					"TestCounter": {counter: 2},
				},
			},
			args: args{
				mtype: "counter",
				mname: "TestCounter",
			},
			body:       `{"id":"TestCounter","type":"counter","delta":2}`,
			statuscode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storedData{
				data: tt.fields.data,
			}
			got, got1 := s.GetStoredDataByParamToJSON(tt.args.mtype, tt.args.mname)
			if !reflect.DeepEqual(string(got), tt.body) {
				t.Errorf("GetStoredDataByParamToJSON() got = %v, err %v", string(got), tt.body)
			}
			if got1 != tt.statuscode {
				t.Errorf("GetStoredDataByParamToJSON() got1 = %v, err %v", got1, tt.statuscode)
			}
		})
	}
}
