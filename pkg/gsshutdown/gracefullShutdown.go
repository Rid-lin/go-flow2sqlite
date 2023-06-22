// package main

// import (
// 	"log"
// 	"os"
// 	"os/signal"
// 	"syscall"
// )

// func getExitSignalsChannel() chan os.Signal {

// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c,
// 		// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
// 		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
// 		syscall.SIGINT,  // Ctrl+C
// 		syscall.SIGQUIT, // Ctrl-\
// 		// syscall.SIGKILL, // "always fatal", "SIGKILL and SIGSTOP may not be caught by a program"
// 		// syscall.SIGHUP, // "terminal is disconnected"
// 	)
// 	return c

// }

// func (t *Transport) Exit(cfg *Config) {
// 	<-t.exitChan
// 	t.fileDestination.Close()
// 	t.conn.Close()
// 	log.Println("Shutting down")
// 	os.Exit(0)

// }

// func getNewLogSignalsChannel() chan os.Signal {

// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c,
// 		// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
// 		syscall.SIGHUP, // "terminal is disconnected"
// 	)
// 	return c

// }

// newLogChan := getNewLogSignalsChannel()

// go func() {
// 	<-newLogChan
// 	log.Println("Received a signal from logrotate, close the file.")
// 	writer.Flush()
// 	filetDestination.Close()
// 	log.Println("Opening a new file.")
// 	time.Sleep(2 * time.Second)
// 	writer = openOutputDevice(cfg.NameFileToLog)

// }()

package gsshutdown

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type GSS struct {
	ExitChan chan os.Signal
	LogChan  chan os.Signal
}

func NewGSS(fnExitfunc func(ve ...interface{}), fnRotate func(vr ...interface{})) *GSS {
	var gss GSS
	ExitChan := getExitSignalsChannel()
	gss.ExitChan = ExitChan
	LogChan := getNewLogSignalsChannel()
	gss.LogChan = LogChan
	go gss.Exit(fnExitfunc)
	go gss.GetSIGHUP(fnRotate)
	return &gss
}

func getExitSignalsChannel() chan os.Signal {

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		// syscall.SIGKILL, // "always fatal", "SIGKILL and SIGSTOP may not be caught by a program"
		// syscall.SIGHUP, // "terminal is disconnected"
	)
	return c

}

func (gss *GSS) Exit(fnExit func(ve ...interface{})) {
	<-gss.ExitChan
	fnExit(nil)
	log.Println("Shutting down")
	// time.Sleep(5 * time.Second)
	os.Exit(0)

}

func (gss *GSS) GetSIGHUP(fnRotate func(vr ...interface{})) {
	<-gss.LogChan
	fnRotate(nil)
}

func getNewLogSignalsChannel() chan os.Signal {

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
		syscall.SIGHUP, // "terminal is disconnected"
	)
	return c

}
