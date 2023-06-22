package main

import (
	"bytes"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	// "git.vegner.org/vsvegner/go-flow2sqlite/cmd/app/internal/config"
	"go-flow2sqlite/cmd/app/internal/config"
	"go-flow2sqlite/cmd/app/internal/services/postprocessing"
	decodeflow "go-flow2sqlite/cmd/app/internal/services/udp"
	"go-flow2sqlite/pkg/gsshutdown"
)

func main() {
	var err error

	cfg := config.NewConfig()

	/*Creating a channel to intercept the program end signal*/

	gss := gsshutdown.NewGSS(
		Exit(nil),
		GetSIGHUP(cfg),
	)

	t := NewApp(cfg)

	go gss.Exit(nil)

	// Endless file parsing loop
	go func(cfg *config.Config) {
		t.s.runOnce(cfg)
		for {
			<-c.timerUpdatedevice.C
			c.runOnce(cfg)
		}
	}(cfg)

	go func(cfg *config.Config, t *Transport) {
		if err := http.ListenAndServe(cfg.BindAddr, t.router); err != nil {
			logrus.Fatal(err)
		}
	}(cfg, t)

	/* Create output pipe */
	outputChannel := make(chan decodeflow.DecodedRecord, 100)
	storeChannel := make(chan postprocessing.BMes, 100)
	go postprocessing.PreparingForStore(storeChannel, outputChannel, cfg.SubNets, cfg.IgnorList, t.ds)

	/* Start listening on the specified port */
	logrus.Infof("Start listening to NetFlow stream on %v", cfg.FlowAddr)
	addr, err := net.ResolveUDPAddr("udp", cfg.FlowAddr)
	if err != nil {
		logrus.Fatalf("Error: %v\n", err)
	}

	for {
		t.conn, err = net.ListenUDP("udp", addr)
		if err != nil {
			logrus.Errorln(err, "Sleeping 5 second")
			time.Sleep(5 * time.Second)
		} else {
			err = t.conn.SetReadBuffer(cfg.ReceiveBufferSizeBytes)
			if err != nil {
				logrus.Errorln(err, "Sleeping 2 second")
				time.Sleep(2 * time.Second)
			} else {
				/* Infinite-loop for reading packets */
				for {
					buf := make([]byte, 4096)
					rlen, remote, err := t.conn.ReadFromUDP(buf)

					if err != nil {
						logrus.Errorf("Error: %v\n", err)
					} else {

						stream := bytes.NewBuffer(buf[:rlen])

						go decodeflow.HandlePacket(stream, remote, outputChannel)
					}
				}
			}
		}

	}
}
