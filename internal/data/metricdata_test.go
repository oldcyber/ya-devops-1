package data

//func TestSetRandomValue(t *testing.T) {
//	tests := []struct {
//		name string
//		want float64
//	}{
//		{
//			name: "test1",
//			want:
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := SetRandomValue(); got != tt.want {
//				t.Errorf("SetRandomValue() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func Test_metricStore_AddMetrics(t *testing.T) {
//	type fields struct {
//		data map[string]gauge
//		mtx  sync.RWMutex
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ms := &metricStore{
//				data: tt.fields.data,
//				mtx:  tt.fields.mtx,
//			}
//			ms.AddMetrics()
//		})
//	}
//}
//
//func Test_metricStore_GetMetrics(t *testing.T) {
//	type fields struct {
//		data map[string]gauge
//		mtx  sync.RWMutex
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   map[string]gauge
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ms := &metricStore{
//				data: tt.fields.data,
//				mtx:  tt.fields.mtx,
//			}
//			if got := ms.GetMetrics(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetMetrics() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
