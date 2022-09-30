package main

import (
	"github.com/quantstop/quantstopterminal/internal/engine"
	"github.com/quantstop/quantstopterminal/internal/log"
)

var (
	// these will be replaced by the release process using ldflags
	version = "v0.0.0"
	commit  = "0000000"
	date    = "0001-01-01T00:00:00Z"
)

func main() {

	log.Debugf(log.Global, "Creating Engine v%s %s (%s)\n", version, commit, date)
	err, bot := engine.NewEngine(version, commit, date)
	if err != nil {
		log.Fatalf(log.Global, "Creating Engine ... Error: %s\n", err)
	}
	log.Debugln(log.Global, "Engine created.")

	log.Debugln(log.Global, "Starting Engine ...")
	bot.Start()

	log.Debugln(log.Global, "Engine started.")
	bot.WaitForInterrupt()
}
