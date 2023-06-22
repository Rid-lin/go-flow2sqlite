package updatergomtc

import (
	"encoding/json"
	"go-flow2sqlite/cmd/app/internal/models"
	"io"
	"net/http"
	"time"
)

type DeviceRepository struct {
	store             *Store
	devices           []*models.DeviceType
	timerUpdatedevice *time.Timer
}

func (r *DeviceRepository) GetAll(addr string, retry uint) ([]*models.DeviceType, error) {
	ds := []*models.DeviceType{}

	return ds, nil
}

func JSONClient(server, uri string) ([]models.DeviceOfMikrotik, error) {
	url := server + uri

	spaceClient := http.Client{
		Timeout: time.Second * 10, // Timeout after 10 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		return nil, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}
	d := []models.DeviceOfMikrotik{}
	jsonErr := json.Unmarshal(body, &d)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return d, nil
}

func (du *DeviceRepository) RunOnce(addr, interval string, retry uint) {
	devices, _ := du.GetAll(addr, retry)
	du.devices = devices
	du.SetTimerUpdateDevice(interval)

}

func (du *DeviceRepository) SetTimerUpdateDevice(IntervalStr string) {
	interval, err := time.ParseDuration(IntervalStr)
	if err != nil {
		du.timerUpdatedevice = time.NewTimer(15 * time.Minute)
	} else {
		du.timerUpdatedevice = time.NewTimer(interval)
	}
}

func (du *DeviceRepository) GetTimer() *time.Timer {
	return du.timerUpdatedevice
}
