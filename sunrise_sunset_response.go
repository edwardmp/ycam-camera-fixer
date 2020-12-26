package main

import (
	"fmt"
	"time"
)

type sunriseSunsetResponse struct {
	sunriseSunsetData `json:"results"`
	Status            string `json:"status"`
	LastFetch         time.Time
}

type sunriseSunsetData struct {
	Sunrise time.Time `json:"sunrise"`
	Sunset  time.Time `json:"sunset"`
}

func (s sunriseSunsetResponse) String() string {
	return fmt.Sprintf(
		"[Status: %s], [LastFetch: %s], [Sunrise: %s], [Sunset: %s]",
		s.Status,
		s.LastFetch,
		s.Sunrise,
		s.Sunset)
}

func (s sunriseSunsetResponse) DataIsOutdatedComparedTo(comparisonTime time.Time) bool {
	y1, m1, d1 := s.LastFetch.Date()
	y2, m2, d2 := comparisonTime.Date()

	return s.Status != "OK" || !(y1 == y2 && m1 == m2 && d1 == d2)
}
