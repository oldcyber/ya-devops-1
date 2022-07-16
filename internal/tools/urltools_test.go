package tools

import (
	"reflect"
	"testing"
)

func TestGetURL(t *testing.T) {
	type args struct {
		url string
		h   string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Simple URL string",
			args: args{
				url: "/update/gauge/Alloc/1234345",
				h:   "update",
			},
			want: []string{"gauge", "Alloc", "1234345"},
		},
		{
			name: "short URL string without one argument",
			args: args{
				url: "/update/gauge/1234345",
				h:   "update",
			},
			want: []string{"gauge", "1234345"},
		},
		{
			name: "Short URL string without update",
			args: args{
				url: "gauge/Alloc/1234345",
				h:   "update",
			},
			want: []string{"gauge", "Alloc", "1234345"},
		},
		{
			name: "wrong long URL string",
			args: args{
				url: "/v1/update/gauge/Alloc/1234345",
				h:   "update",
			},
			want: []string{"v1", "update", "gauge", "Alloc", "1234345"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetURL(tt.args.url, tt.args.h); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
