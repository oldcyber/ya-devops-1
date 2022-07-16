package server

import (
	"testing"
)

func Test_storeData1(t *testing.T) {
	type args struct {
		res []string
	}
	type Want struct {
		err bool
		an  int
	}
	tests := []struct {
		name string
		args args
		want Want
	}{
		{
			name: "simple test",
			args: args{
				res: []string{"gauge", "Alloc", "100"},
			},
			want: Want{
				err: true,
				an:  200,
			},
		},
		{
			name: "not simple test",
			args: args{
				res: []string{"gauge", "Test", "100"},
			},
			want: Want{
				err: false,
				an:  400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := storeData(tt.args.res)
			if (got != tt.want.err) || (got1 != tt.want.an) {
				t.Errorf("storeData() = %v, %v, want %v, %v", got, got1, tt.want.err, tt.want.an)
			}
			//if got, err := storeData(tt.args.res); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("storeData() = %v, want %v, err = %v", got, tt.want, err)
			//}
		})
	}
}
