package main

import (
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/engine"
	"github.com/quantstop/quantstopterminal/internal/log"
	"os"
	"strings"
)

var (
	// these will be replaced by goreleaser
	version = "v0.0.0"
	commit  = "0000000"
	date    = "0001-01-01T00:00:00Z"
)

func main() {

	// parse goreleaser info
	if len(os.Args) >= 2 && "version" == strings.TrimPrefix(os.Args[1], "") {
		fmt.Printf("YOUR_CLI_NAME v%s %s (%s)\n", version, commit[:7], date)
	}

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
