package engine

import (
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/database"
	"github.com/quantstop/quantstopterminal/internal/engine/config"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/internal/system"
	"github.com/quantstop/quantstopterminal/internal/webserver"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

var engineLock = &sync.Mutex{}
var engineInstance *Engine

type Engine struct {
	*config.Config  // Config for engine, and all services
	*ServiceManager // Registry of all services
	sync.WaitGroup  // Service WaitGroup
}

func NewEngine(version, commit, date string) (error, *Engine) {
	if engineInstance == nil {
		engineLock.Lock()
		defer engineLock.Unlock()
		if engineInstance == nil {

			var err error

			engineInstance = &Engine{}

			log.Debugln(log.Global, "Creating Config ...")
			if engineInstance.Config, err = config.NewConfig(); err != nil {
				return err, nil
			}

			log.Debugln(log.Global, "Verifying Config ...")
			if err = engineInstance.Config.CheckConfig(); err != nil {
				log.Errorf(log.Global, "Error checking config: %s\n", err)
			}

			log.Debugln(log.Global, "Creating ServiceManager ...")
			engineInstance.ServiceManager = NewServiceManager()

			log.Debugln(log.Global, "Registering Services ...")
			if err = engineInstance.registerServices(); err != nil {
				return err, nil
			}

			// Set the bot version
			//bot.Version = version

			// Set the max processors for go
			if err = system.AdjustGoMaxProcs(engineInstance.Config.GoMaxProcessors); err != nil {
				return fmt.Errorf("unable to adjust runtime GOMAXPROCS value. Err: %s", err), nil
			}

			// Print banner and version
			//log.Infof(log.Global, "\n"+banner.GetRandomBanner()+"\n"+version.GetVersionString(false))

			// Print info
			log.Debugln(log.Global, "Logger initialized.")
			log.Debugf(log.Global, "Using config dir: %s\n", engineInstance.Config.ConfigDir)
			if strings.Contains(engineInstance.Config.Logger.Output, "file") {
				log.Debugf(log.Global, "Using log file: %s\n",
					filepath.Join(log.FilePath, engineInstance.Config.Logger.LoggerFileConfig.FileName))
			}

		} else {
			log.Debug(log.Global, "Engine instance already created.")
		}
	} else {
		log.Debug(log.Global, "Engine instance already created.")
	}

	return nil, engineInstance
}

func (bot *Engine) registerServices() (err error) {

	if err = bot.registerDatabaseService(); err != nil {
		return err
	}

	if err = bot.registerWebserverService(); err != nil {
		return err
	}

	return nil
}

func (bot *Engine) registerDatabaseService() error {

	// Create Database
	db, err := database.NewDatabase(bot.Config.Database)
	if err != nil {
		return err
	}

	// Register Database
	if err = bot.RegisterService(db); err != nil {
		return err
	}

	return nil
}

func (bot *Engine) registerWebserverService() error {

	// Fetch dependencies
	var db *database.Database
	if err := bot.FetchService(&db); err != nil {
		log.Error(log.Global, err)
	}

	// Create Webserver
	ws, err := webserver.NewWebserver(bot.Config.Webserver, db)
	if err != nil {
		return err
	}

	// Register Webserver
	if err = bot.RegisterService(ws); err != nil {
		return err
	}

	return nil
}

func (bot *Engine) Start() {
	bot.StartAll(&bot.WaitGroup)
}

func (bot *Engine) Stop() {
	bot.StopAll()
	bot.WaitGroup.Wait()
}

func (bot *Engine) Restart() error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	args := os.Args
	env := os.Environ()
	// Windows does not support exec syscall.
	if runtime.GOOS == "windows" {
		cmd := exec.Command(self, args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = env
		err := cmd.Run()
		if err == nil {
			os.Exit(0)
		}
		return err
	}
	return syscall.Exec(self, args, env)
}

func (bot *Engine) WaitForInterrupt() {
	// Wait for system interrupt to stop the bot
	interrupt := system.WaitForInterrupt()
	log.Infof(log.Global, "Captured %v, shutdown requested.", interrupt)
	bot.Stop()
	log.Infof(log.Global, "Exiting.")
}
