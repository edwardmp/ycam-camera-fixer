# Y-cam Camera Fixer

Uses an API to determine sunrise/sunsit times, and changes the night vision settings of the camera accordingly, to account for the broken camera light sensor.

## To build

```bash
docker build . -t ycam-camera-fixer
```

## Enviroment variables to set

```bash
CAMERA_IP
AUTH_USERNAME
AUTH_PASSWORD
CAMERA_LOCATION_LATITUDE
CAMERA_LOCATION_LONGITUDE
```

