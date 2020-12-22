package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type cameraFixer struct {
	sunriseSunsetResponse sunriseSunsetResponse
	sunriseChangeDone     bool
	sunsetChangeDone      bool
	config                config
}

func (c *cameraFixer) run() {
	defer time.Sleep(1 * time.Second)

	now := time.Now()
	// check if data was fetched last yesterday
	if c.sunriseSunsetResponse.DataIsOutdatedComparedTo(now) {
		err := c.updateSunriseSunsetResponse()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		return
	}

	if !c.sunriseChangeDone && now.After(c.sunriseSunsetResponse.Sunrise) && now.Before(c.sunriseSunsetResponse.Sunset) {
		err := c.changeCameraSettingsForDay()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		c.sunriseChangeDone = true
		return
	}

	if !c.sunsetChangeDone && now.After(c.sunriseSunsetResponse.Sunset) {
		err := c.changeCameraSettingsForNight()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		c.sunsetChangeDone = true
	}
}

func (c *cameraFixer) updateSunriseSunsetResponse() error {
	// Create request
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02")
	urlFormat := "https://api.sunrise-sunset.org/json?lat=%f&lng=%f&date=%s&formatted=0"
	requestUrl := fmt.Sprintf(urlFormat, c.config.CameraLocationLatitude, c.config.CameraLocationLongitude, currentDate)

	resp, err := http.Get(requestUrl)
	if err != nil {
		return fmt.Errorf("could not get sunrise/sunset data: %s", err)
	}

	// Read Response Body
	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	sunriseSunsetResponse := sunriseSunsetResponse{}
	jsonErr := json.Unmarshal(respBody, &sunriseSunsetResponse)
	if jsonErr != nil {
		return fmt.Errorf("sunrise/sunset API response could not be unmarshalled: %s", jsonErr)
	}

	if sunriseSunsetResponse.Status != "OK" {
		return fmt.Errorf("sunrise/sunset API status is not OK: %s", sunriseSunsetResponse.Status)
	}

	// convert the configured time zone from UTC
	sunriseSunsetResponse.Sunset = c.timeInConfiguredTimezone(sunriseSunsetResponse.Sunset)
	sunriseSunsetResponse.Sunrise = c.timeInConfiguredTimezone(sunriseSunsetResponse.Sunrise)

	sunriseSunsetResponse.LastFetch = time.Now()
	log.Print("Sunset/sunrise times fetched: ", sunriseSunsetResponse)

	c.sunriseSunsetResponse = sunriseSunsetResponse

	// reset
	c.sunriseChangeDone = false
	c.sunsetChangeDone = false

	return nil
}

func (c *cameraFixer) changeCameraSettingsForDay() error {
	params := url.Values{}
	params.Set("IRLED", "off")
	params.Set("BWMODE", "off")
	params.Set("MLMODE", "off")
	params.Set("IRCUT", "off")

	log.Print("Changing camera settings for day..")
	return c.changeCameraSettings(params)
}

func (c *cameraFixer) changeCameraSettingsForNight() error {
	params := url.Values{}
	params.Set("IRLED", "on")
	params.Set("BWMODE", "on")
	params.Set("MLMODE", "on")
	params.Set("IRCUT", "on")

	log.Print("Changing camera settings for night..")
	return c.changeCameraSettings(params)
}

func (c *cameraFixer) changeCameraSettings(params url.Values) error {
	body := bytes.NewBufferString(params.Encode())

	client := &http.Client{}

	requestUrl := fmt.Sprintf("http://%s/form/nvctlApply", c.config.CameraIP)
	req, err := http.NewRequest("POST", requestUrl, body)
	if err != nil {
		return fmt.Errorf("could not create request: %s", err)
	}

	req.SetBasicAuth(c.config.AuthUsername, c.config.AuthPassword)

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

func (c *cameraFixer) timeInConfiguredTimezone(timeToConvert time.Time) time.Time {
	loc, err := time.LoadLocation(c.config.TimeZone)
	if err != nil {
		log.Fatal(err)
	}

	return timeToConvert.In(loc)
}
