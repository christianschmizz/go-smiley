package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"

	ssp "github.com/christianschmizz/go-smiley"
	"github.com/rs/zerolog"
)

var (
	portName string
)

func main() {
	flag.StringVar(&portName, "port", "/dev/cu.usbmodem1422401", "port name")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	if portName == "" {
		log.Fatal().Msg("please provide a port to use")
	}

	nvx, err := ssp.DialNVx(portName)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to device")
	}
	defer nvx.Close()

	if err := nvx.Init(); err != nil {
		log.Fatal().Err(err).Msg("device initialization failed")
	}

	nvx.Poll(c)
}
