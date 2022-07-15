package agent

import "testing"

func Test_sendCounterMetrics(t *testing.T) {
	type args struct {
		c counter
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendCounterMetrics(tt.args.c)
		})
	}
}

func Test_sendGaugeMetrics(t *testing.T) {
	type args struct {
		m map[string]gauge
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendGaugeMetrics(tt.args.m)
		})
	}
}
