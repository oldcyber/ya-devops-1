package server

import (
	"testing"
)

func Test_storeData(t *testing.T) {
	type args struct {
		res []string
	}

	type Data struct {
		Mtype string
		Name  string
		Val   storedType
	}
	type want struct {
		Data Data
		ok   bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "simple test",
			args: args{
				res: []string{"gauge", "Alloc", "100"},
			},
			want: want{
				Data: Data{
					Mtype: "gauge",
					Name:  "Alloc",
					Val:   storedType{gauge: 100},
				},
				ok: true,
			},
		},
		{
			name: "not simple test",
			args: args{
				res: []string{"gauge", "Alloc", "100", "вентилятор"},
			},
			want: want{
				Data: Data{
					Mtype: "gauge",
					Name:  "Alloc",
					Val:   storedType{gauge: 100},
				},
				ok: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StoredData = make(map[int]*SData)
			_, ok := storeData(tt.args.res)
			//if !reflect.DeepEqual(got[0], tt.want.Data) {
			// t.Errorf("storeData() got = %v, want %v", got[0], tt.want.Data)
			//}
			if ok != tt.want.ok {
				t.Errorf("storeData() ok = %v, want %v", ok, tt.want.ok)
			}
		})
	}
}
