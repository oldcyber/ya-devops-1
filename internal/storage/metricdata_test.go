package storage

import (
	"testing"
)

func TestMetricStore_AddMetrics(t *testing.T) {
	type fields struct {
		data map[string]gauge
	}
	tests := []struct {
		name   string
		fields fields
		arg    chan float64
	}{
		{
			name: "test",
			arg:  make(chan float64, 3),
			fields: fields{
				data: map[string]gauge{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MetricStore{
				data: tt.fields.data,
			}
			go ms.GetNewMetrics(tt.arg)
			ms.AddMetrics(tt.arg)
		})
	}
}
