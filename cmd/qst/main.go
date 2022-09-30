package main

import (
	"github.com/quantstop/quantstopterminal/internal/engine"
	"github.com/quantstop/quantstopterminal/internal/log"
)

func main() {

	log.Debugln(log.Global, "Creating Engine ...")
	err, bot := engine.NewEngine()
	if err != nil {
		log.Fatalf(log.Global, "Creating Engine ... Error: %s\n", err)
	}
	log.Debugln(log.Global, "Creating Engine ... Success.")

	log.Debugln(log.Global, "Starting Engine ...")
	bot.Start()

	log.Debugln(log.Global, "Starting Engine ... Success.")
	bot.WaitForInterrupt()
}
