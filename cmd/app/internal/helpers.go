package main

import (
	"database/sql"
	"go-flow2sqlite/cmd/app/internal/config"
	"go-flow2sqlite/cmd/app/internal/models"
	serverHTTP "go-flow2sqlite/cmd/app/internal/servers/http"
	"go-flow2sqlite/cmd/app/internal/services/postprocessing"
	"net"
	"os"
	"sync"
	"time"
)

func Exit(ve ...interface{}) func(ve ...interface{}) {
	return func(ve ...interface{}) {
	}
}

func GetSIGHUP(vr ...interface{}) func(vr ...interface{}) {
	return func(vr ...interface{}) {
	}
}

func (s *Store) writeMessageToStore(inputChannel chan postprocessing.BMes) {
	for bm := range inputChannel {
		s.SaveLineOfLog(bm)
		// ctx := context.Background()
		// s.DB.ExecContext(ctx, "INSERT INTO raw_log VALUES(?,?,?,?,?,?,?,?,?,?)",
		// 	bm.Time,
		// 	bm.Delay,
		// 	bm.IPSrc,
		// 	bm.Protocol,
		// 	bm.SizeOfByte,
		// 	bm.IPDst,
		// 	bm.PortDst,
		// 	bm.MacDst,
		// 	bm.RouterIP,
		// 	bm.PortSRC,
		// )

	}

}

type Transport struct {
	DB                *sql.DB
	ds                *models.Devices
	s                 *serverHTTP.Server
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

func NewApp(cfg *config.Config) *Transport {
	var err error

	Location := time.UTC
	db, err := sql.Open("sqlite", cfg.PathToBD)
	if err != nil {
		os.Exit(2)
	}

	t := &Transport{
		DB:          db,
		renewOneMac: make(chan string, 100),
		s:           serverHTTP.NewServer(cfg),
		Location:    Location,
		Interval:    cfg.Interval,
		GomtcAddr:   cfg.GomtcAddr,
	}
	return t
}
