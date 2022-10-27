package tools

import (
	"reflect"
	"testing"
)

func TestSetURL(t *testing.T) {
	type args struct {
		k string
		v string
		t string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Simple gauge URL string",
			args: args{
				k: "Alloc",
				v: "1234345",
				t: "gauge",
			},
			want: "http://127.0.0.1:8080/update/gauge/Alloc/1234345",
		},
		{
			name: "Simple counter URL string",
			args: args{
				k: "Alloc",
				v: "1234345",
				t: "counter",
			},
			want: "http://127.0.0.1:8080/update/counter/Alloc/1234345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetURL(tt.args.k, tt.args.v, tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
