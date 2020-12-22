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
	return s.Status != "OK" || !s.LastFetch.Truncate(time.Hour*24).Equal(comparisonTime.Truncate(time.Hour*24))
}
