package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

type ResponseType struct {
	IP       string `JSON:"IP"`
	Mac      string `JSON:"Mac"`
	HostName string `JSON:"Hostname"`
	Comments string `JSON:"Comment"`
}

type Transport struct {
	DB *sql.DB
	// cache             *timedmap.TimedMap
	ipToMac           map[string]LineOfData
	router            *mux.Router
	Location          *time.Location
	fileDestination   *os.File
	conn              *net.UDPConn
	timerUpdatedevice *time.Timer
	renewOneMac       chan string
	exitChan          chan os.Signal
	Interval          string
	GomtcAddr         string
	sync.RWMutex
}

func NewTransport(cfg *Config) *Transport {
	var err error

	fileDestination, err = os.OpenFile(cfg.PathToBD, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fileDestination.Close()
		logrus.Fatalf("Error, the '%v' file could not be created (there are not enough premissions or it is busy with another program): %v", cfg.PathToBD, err)
	}

	Location := time.UTC
	db, err := sql.Open("sqlite", cfg.PathToBD)
	if err != nil {
		os.Exit(2)
	}

	t := &Transport{
		DB: db,
		// cache:           timedmap.New(600 * time.Second),
		// TODO следующую строку можно заменить на sqlite в памяти, всёравно её уже используем. Подумай над этим.
		ipToMac:         make(map[string]LineOfData),
		renewOneMac:     make(chan string, 100),
		router:          mux.NewRouter(),
		Location:        Location,
		exitChan:        getExitSignalsChannel(),
		Interval:        cfg.Interval,
		fileDestination: fileDestination,
		GomtcAddr:       cfg.GomtcAddr,
	}
	t.configureRouter()
	return t
}

// func (data *Transport) GetInfo(ip, time string) ResponseType {
// 	var response ResponseType
// 	d, ok := data.cache.GetValue(ip).(string)
// 	if !ok {
// 		data.timerUpdatedevice.Stop()
// 		data.setTimerUpdateDevice(data.Interval)

// 	}

// 	response.Mac = d
// 	if response.Mac == "" {
// 		response.Mac = ip
// 	}
// 	if response.HostName == "" {
// 		response.HostName = "-"
// 	}

// 	return response
// }

func (data *Transport) GetInfo(ip, time string) ResponseType {
	var response ResponseType
	data.RLock()
	ipStruct, ok := data.ipToMac[ip]
	data.RUnlock()
	if ok {
		response.Mac = ipStruct.Mac
		if response.Mac == "" {
			response.Mac = ipStruct.ActiveMacAddress
		}
		response.HostName = ipStruct.HostName
	} else {
		data.timerUpdatedevice.Stop()
		data.setTimerUpdateDevice(data.Interval)
	}
	if response.Mac == "" {
		response.Mac = ip
	}
	if response.HostName == "" {
		response.HostName = "-"
	}

	return response
}

func getDataOverApi(
	// qh, qd, qm uint64,
	addr string) map[string]LineOfData {
	lineOfData := LineOfData{}
	ipToMac := map[string]LineOfData{}
	// arrDevices := []Device{}
	arrDevices, err := JSONClient(addr, "/api/v1/devices")
	if err != nil {
		logrus.Error(err)
		return ipToMac
	}
	for _, value := range arrDevices {
		lineOfData.Device = value
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

func JSONClient(server, uri string) ([]Device, error) {
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

func (t *Transport) getDevices() map[string]LineOfData {
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
