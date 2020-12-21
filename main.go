package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"net/http"
	"os"
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

func (s sunriseSunsetResponse) DataIsOutdated() bool {
	now := time.Now()
	return !s.LastFetch.Truncate(time.Hour * 24).Equal(now.Truncate(time.Hour * 24))
}

func main() {
	var sunriseSunsetResponse *sunriseSunsetResponse = nil
	sunriseChangeDone := false
	sunsetChangeDone := false

	for {
		now := time.Now()
		// check if data was fetched last yesterday
		if sunriseSunsetResponse == nil || sunriseSunsetResponse.DataIsOutdated() {
			response, err := getSunsetSunriseTimes()
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			sunriseSunsetResponse = response

			sunriseChangeDone = false
			sunsetChangeDone = false
		} else if !sunriseChangeDone && now.After(sunriseSunsetResponse.Sunrise) && now.Before(sunriseSunsetResponse.Sunset) {
			err := changeCameraSettingsForDay()
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			sunriseChangeDone = true
		} else if !sunsetChangeDone && now.After(sunriseSunsetResponse.Sunset) {
			err := changeCameraSettingsForNight()
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			sunsetChangeDone = true
		}
	}
}

func getSunsetSunriseTimes() (*sunriseSunsetResponse, error) {
	// Create request
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02")
	requestUrl := fmt.Sprintf("https://api.sunrise-sunset.org/json?lat=52.142501&lng=4.396298&date=%s&formatted=0", currentDate)

	resp, err := http.Get(requestUrl)
	if err != nil {
		return nil, fmt.Errorf("could not get sunrise/sunset data: %s", err)
	}

	// Read Response Body
	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	sunriseSunsetResponse := sunriseSunsetResponse{}
	jsonErr := json.Unmarshal(respBody, &sunriseSunsetResponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	sunriseSunsetResponse.LastFetch = time.Now()
	log.Print("Sunset/sunrise times fetched: ", sunriseSunsetResponse)

	return &sunriseSunsetResponse, nil
}

func changeCameraSettingsForDay() error {
	params := url.Values{}
	params.Set("IRLED", "off")
	params.Set("BWMODE", "off")
	params.Set("MLMODE", "off")
	params.Set("IRCUT", "off")

	log.Print("Changing camera settings for day..")
	return changeCameraSettings(params)
}

func changeCameraSettingsForNight() error {
	params := url.Values{}
	params.Set("IRLED", "on")
	params.Set("BWMODE", "on")
	params.Set("MLMODE", "on")
	params.Set("IRCUT", "on")

	log.Print("Changing camera settings for night..")
	return changeCameraSettings(params)
}

func changeCameraSettings(params url.Values) error {
	body := bytes.NewBufferString(params.Encode())

	client := &http.Client{}

	cameraIp := os.Getenv("CAMERA_IP")
	if cameraIp == "" {
		log.Fatal("CAMERA_IP env variable is not set")
	}

	requestUrl := fmt.Sprintf("http://%s/form/nvctlApply", cameraIp)
	req, err := http.NewRequest("POST", requestUrl, body)
	if err != nil {
		return fmt.Errorf("could not create request: %s", err)
	}

	authUsername := os.Getenv("AUTH_USERNAME")
	if authUsername == "" {
		log.Fatal("AUTH_USERNAME env variable is not set")
	}

	authPassword := os.Getenv("AUTH_PASSWORD")
	if authPassword == "" {
		log.Fatal("authPassword env variable is not set")
	}
	req.SetBasicAuth(authUsername, authPassword)

	// Headers
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not change camera settings: %s", err)
	}

	// Display Results
	log.Print("Camera response status: ", resp.Status)

	return nil
}
