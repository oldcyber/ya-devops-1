package storage

import (
	"testing"

	"github.com/oldcyber/ya-devops-1/internal/env"
)

func TestStoredData_AddNewItem(t *testing.T) {
	type fields struct {
		StoredType StoredType
		s          *StoredMem
	}
	type args struct {
		res []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  int
	}{
		{
			name: "test1",
			fields: fields{
				StoredType: StoredType{
					gauge:   0,
					counter: 0,
					stype:   "",
				},
				s: &StoredMem{
					data: map[string]StoredType{},
				},
			},
			args: args{
				res: []string{env.MetricGaugeType, "test", "1"},
			},
			want:  true,
			want1: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := &StoredData{
				StoredType: tt.fields.StoredType,
				s:          tt.fields.s,
			}
			got, got1 := sd.AddNewItem(tt.args.res)
			if got != tt.want {
				t.Errorf("AddNewItem() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AddNewItem() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
