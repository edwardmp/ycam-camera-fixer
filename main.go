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
	Status string `json:"status"`
	LastFetch time.Time
}

type sunriseSunsetData struct {
	Sunrise time.Time `json:"sunrise"`
	Sunset time.Time `json:"sunset"`
}

func main() {
	var sunriseSunsetResponse *sunriseSunsetResponse = nil
	sunriseChangeDone := false
	sunsetChangeDone := false

	for {
		now := time.Now()
		// check if data was fetched last yesterday
		if (sunriseSunsetResponse == nil || !sunriseSunsetResponse.LastFetch.Truncate(time.Hour * 24).Equal(now.Truncate(time.Hour * 24))) {
			response := getSunsetSunriseTimes()
			sunriseSunsetResponse = &response

			sunriseChangeDone = false
			sunsetChangeDone = false
		}

		if (!sunriseChangeDone && now.After(sunriseSunsetResponse.Sunrise) && now.Before(sunriseSunsetResponse.Sunset)) {
			changeCameraSettingsForDay()
			sunriseChangeDone = true
		} else if (!sunsetChangeDone && now.After(sunriseSunsetResponse.Sunset)) {
			changeCameraSettingsForNight()
			sunsetChangeDone = true
		}
	}
}

func getSunsetSunriseTimes() sunriseSunsetResponse {
	// Create request
	currentTime := time.Now()
	url := fmt.Sprintf("https://api.sunrise-sunset.org/json?lat=52.142501&lng=4.396298&date=%s&formatted=0", currentTime.Format("2006-01-02"))

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Failure : ", err)
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

	log.Print("Sunset/sunrise times fetched: ", sunriseSunsetResponse.Status)

	return sunriseSunsetResponse
}

func changeCameraSettingsForDay() {
	params := url.Values{}
	params.Set("IRLED", "off")
	params.Set("BWMODE", "off")
	params.Set("MLMODE", "off")
	params.Set("IRCUT", "off")

	log.Print("Changing camera settings for day..")
	changeCameraSettings(params)
}

func changeCameraSettingsForNight() {
	params := url.Values{}
	params.Set("IRLED", "on")
	params.Set("BWMODE", "on")
	params.Set("MLMODE", "on")
	params.Set("IRCUT", "on")

	log.Print("Changing camera settings for night..")
	changeCameraSettings(params)
}

func changeCameraSettings(params url.Values) {
	body := bytes.NewBufferString(params.Encode())

	// Create client
	client := &http.Client{}

	// Create request
	cameraIp := os.Getenv("CAMERA_IP")
	if cameraIp == "" {
		log.Fatal("CAMERA_IP env variable is not set")
	}

	url := fmt.Sprintf("http://%s/form/nvctlApply", cameraIp)
	req, err := http.NewRequest("POST", url, body)

	// Headers
	basicAuthEncoded := os.Getenv("BASIC_AUTH")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", basicAuthEncoded))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Display Results
	log.Print("Response Status: ", resp.Status)
}