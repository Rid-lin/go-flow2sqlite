package main

import (
	"go-flow2sqlite/cmd/app/internal/config"
	"go-flow2sqlite/cmd/app/internal/models"
	"go-flow2sqlite/cmd/app/internal/services/updater"
	updatergomtc "go-flow2sqlite/cmd/app/internal/services/updater/gomtc"
	"go-flow2sqlite/cmd/app/internal/store"
	sqlitestore "go-flow2sqlite/cmd/app/internal/store/sqlite"
	"log"
	"net"
	"os"
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

type Transport struct {
	// s                 *serverHTTP.Server
	// renewOneMac       chan string
	ds        *models.Devices
	updater   updater.Store
	store     store.Store
	Location  *time.Location
	conn      *net.UDPConn
	Interval  string
	GomtcAddr string
}

func NewApp(cfg *config.Config) *Transport {
	db, err := sqlitestore.NewDB(cfg.PathToBD)
	if err != nil {
		log.Println("Error open DB: ", err)
		os.Exit(1)
	}
	t := &Transport{
		// s:           serverHTTP.NewServer(cfg),
		// renewOneMac: make(chan string, 100),
		updater:   updatergomtc.NewUpdater(cfg.GomtcAddr),
		store:     sqlitestore.NewStore(db),
		Location:  SetLocation(cfg.Loc),
		Interval:  cfg.Interval,
		GomtcAddr: cfg.GomtcAddr,
	}
	return t
}

func SetLocation(loc string) *time.Location {
	location, err := time.LoadLocation(loc)
	if err != nil {
		location = time.UTC
		println("Error loading Location(", loc, "):", err, " (set default location = UTC)")
	}
	return location
}
