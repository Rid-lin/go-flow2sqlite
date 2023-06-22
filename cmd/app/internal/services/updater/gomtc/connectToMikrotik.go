package updater

import (
	"encoding/json"
	"go-flow2sqlite/cmd/app/internal/config"
	"go-flow2sqlite/cmd/app/internal/models"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func getDataOverApi(
	// qh, qd, qm uint64,
	addr string) map[string]models.LineOfData {
	lineOfData := models.LineOfData{}
	ipToMac := map[string]models.LineOfData{}
	// arrDevices := []Device{}
	arrDevices, err := JSONClient(addr, "/api/v1/devices")
	if err != nil {
		logrus.Error(err)
		return ipToMac
	}
	for _, value := range arrDevices {
		lineOfData.DeviceOfMikrotik = value
		// if value.HourlyQuota == 0 {
		// 	value.HourlyQuota = qh
		// }
		// if value.DailyQuota == 0 {
		// 	value.DailyQuota = qd
		// }
		// if value.MonthlyQuota == 0 {
		// 	value.MonthlyQuota = qm
		// }
		lineOfData.addressLists = strings.Split(lineOfData.AddressLists, ",")
		lineOfData.Timeout = time.Now()
		ipToMac[lineOfData.IP] = lineOfData
	}
	return ipToMac
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
	d := []Device{}
	jsonErr := json.Unmarshal(body, &d)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return d, nil
}

func getDevices() map[string]models.DeviceOfMikrotik {
	t.RLock()
	// qh := t.HourlyQuota
	// qd := t.DailyQuota
	// qm := t.MonthlyQuota
	addr := t.GomtcAddr
	t.RUnlock()
	return getDataOverApi(
		// qh, qd, qm,
		addr)
}

func (t *Transport) runOnce(cfg *config.Config) {
	t.ds = t.getmodels.DeviceOfMikrotik()
	t.setTimerUpdateDevice(cfg.Interval)

}

func (t *Transport) setTimerUpdateDevice(IntervalStr string) {
	t.Lock()
	interval, err := time.ParseDuration(IntervalStr)
	if err != nil {
		t.timerUpdatedevice = time.NewTimer(15 * time.Minute)
	} else {
		t.timerUpdatedevice = time.NewTimer(interval)
	}
	t.Unlock()
}

func (t *Transport) GetDevices() ([]models.LineOfData, error) {
	devices := []models.LineOfData{}
	t.RLock()
	for _, value := range t.ipToMac {
		devices = append(devices, value)
	}
	t.RUnlock()
	return devices, nil
}
