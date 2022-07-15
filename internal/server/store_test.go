package server

import (
	"reflect"
	"testing"
)

func Test_storeData1(t *testing.T) {
	type args struct {
		res []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "simple test",
			args: args{
				res: []string{"gauge", "Alloc", "100"},
			},
			want: true,
		},
		{
			name: "not simple test",
			args: args{
				res: []string{"gauge", "Alloc", "100", "вентилятор"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := storeData(tt.args.res); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("storeData() = %v, want %v, err = %v", got, tt.want, err)
			}
		})
	}
}
