package main

import (
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

type config struct {
	CameraIP                string  `env:"CAMERA_IP,required"`
	AuthUsername            string  `env:"AUTH_USERNAME,required"`
	AuthPassword            string  `env:"AUTH_PASSWORD,required"`
	TimeZone                string  `env:"TZ" envDefault:"UTC"`
	CameraLocationLatitude  float32 `env:"CAMERA_LOCATION_LATITUDE,required"`
	CameraLocationLongitude float32 `env:"CAMERA_LOCATION_LONGITUDE,required"`
	TickIntervalSeconds     int     `env:"TICK_INTERVAL_SECONDS" envDefault:"60"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	cameraFixer := cameraFixer{config: cfg}
	for range time.Tick(time.Duration(cfg.TickIntervalSeconds) * time.Second) {
		cameraFixer.run()
	}
}
