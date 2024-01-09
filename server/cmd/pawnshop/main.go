package main

import (
	"flag"
	"os"
	"os/signal"
	"pawnshop/server/pkg/server"
	"syscall"

	log "github.com/sirupsen/logrus"
)

/*
Runs the pawn shop server.
It accepts two flags: size, which is the size of the inventory, and loglevel, which is the log level.
Defaults to size 2 and log level info.
Also handles graceful shutdown.
*/
func main() {
	invSize := flag.Int("size", 2, "inventory size")
	logLvlStr := flag.String("loglevel", "info", "log level")
	flag.Parse()

	logLvl, err := log.ParseLevel(*logLvlStr)
	if err != nil {
		log.Warnf("Failed to parse log level: %s, using InfoLevel instead", err)
		logLvl = log.InfoLevel
	} else {
		log.Infof("Using log level %s", logLvl)
	}

	log.SetLevel(logLvl)

	server, err := server.New(*invSize)
	if err != nil {
		log.Fatalf("Failed to create server: %s", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Handle a simple graceful shutdown
	go func() {
		<-sigs
		log.Info("Stopping server...")
		err := server.Stop()
		if err != nil {
			log.Errorf("Failed to stop server: %s", err)
		}
	}()

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
