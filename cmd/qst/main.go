package main

import (
	"github.com/quantstop/quantstopterminal/internal/engine"
	"github.com/quantstop/quantstopterminal/internal/log"
)

var (
	// these will be replaced by the build process using ldflags
	version = "0.0.0"
	commit  = "0000000"
	date    = "0001-01-01T00:00:00Z"

	// bot instance
	bot *engine.Engine

	// bot error
	err error
)

func main() {

	// create the engine
	if bot, err = engine.NewEngine(version, commit, date); err != nil {
		log.Fatalf("Error creating engine: %s\n", err)
	}

	// start the engine
	if err = bot.Start(); err != nil {
		log.Fatalf("Error starting engine: %s\n", err)
	}

	// block until interrupt is received
	if err = bot.WaitForInterrupt(); err != nil {
		log.Fatalf("Error waiting for interrupt: %s\n", err)
	}

	// stop the engine
	if err = bot.Stop(); err != nil {
		log.Fatalf("Error stopping engine: %s\n", err)
	}

}
