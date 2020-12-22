package main

import (
	"testing"
	"time"
)

func Test_sunriseSunsetResponse_DataIsOutdatedComparedTo(t *testing.T) {
	referenceTime, _ := time.Parse(
		time.RFC3339,
		"2020-12-22T00:01:41+00:00")

	type fields struct {
		LastFetch time.Time
		Status    string
	}
	type args struct {
		comparisonTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"Similar time should not be considered outdated", fields{referenceTime, "OK"}, args{referenceTime}, false},
		{"Last fetch having same date should not be considered outdated", fields{referenceTime, "OK"}, args{referenceTime.Add(time.Duration(10) * time.Hour)}, false},
		{"Last fetch having another date should be considered outdated", fields{referenceTime.AddDate(0, 0, -1), "OK"}, args{referenceTime}, true},
		{"Not OK status should be considered outdated", fields{referenceTime, "Error"}, args{referenceTime}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := sunriseSunsetResponse{
				Status:    tt.fields.Status,
				LastFetch: tt.fields.LastFetch,
			}
			if got := s.DataIsOutdatedComparedTo(tt.args.comparisonTime); got != tt.want {
				t.Errorf("DataIsOutdatedComparedTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
