package go_mqtt_client

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	debug := flag.Bool("debug", false, "sets log level to debug")
	human := flag.Bool("human", false, "sets log style to human friendly. This is slower.")

	flag.Parse()

	// Set the logger
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if *human {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}
