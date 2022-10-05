package engine

import (
	"errors"
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/engine/banner"
	"github.com/quantstop/quantstopterminal/internal/engine/config"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/internal/system"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

var (
	engineInstance   *Engine
	engineLock       = &sync.Mutex{}
	errEngineCreated = errors.New("engine instance already created")
	errNilEngine     = errors.New("engine instance is nil")
)

type Engine struct {
	*Version        // The engine version information
	*config.Config  // Config for engine, and all services
	*ServiceManager // Registry of all services
}

// NewEngine creates a new engine, and all services
func NewEngine(version, commit, date string) (*Engine, error) {
	var err error

	// singleton pattern enforces that the engine is created only once
	if engineInstance == nil {
		engineLock.Lock()
		defer engineLock.Unlock()
		if engineInstance == nil {

			// parse command line flags
			// ToDo: implement flags?

			// create engine
			engineInstance = &Engine{
				Version: CreateVersion(version, date, commit, false, true, true),
			}

			// create config
			if engineInstance.Config, err = config.NewConfig(); err != nil {
				return nil, err
			}

			// create logger
			if err = log.Initialize(engineInstance.LogConfig); err != nil {
				return nil, err
			}

			// log startup info
			log.Infoln(log.Global, "Creating Engine ...")
			log.Infof(log.Global, "\n"+banner.GetRandomBanner()+"\n"+engineInstance.Version.GetVersionString(false))
			log.Infof(log.Global, "Using config dir: %s\n", engineInstance.ConfigDir)
			if strings.Contains(engineInstance.LogConfig.Output, "file") {
				log.Infof(log.Global, "Using log file: %s\n",
					filepath.Join(log.FilePath, engineInstance.LogConfig.LoggerFileConfig.FileName))
			}

			// log config info
			if err = system.AdjustGoMaxProcs(engineInstance.GoMaxProcessors); err != nil {
				return nil, fmt.Errorf("unable to adjust runtime GOMAXPROCS value. Err: %s", err)
			}

			// create service manager, and all services
			if engineInstance.ServiceManager, err = NewServiceManager(engineInstance); err != nil {
				return nil, err
			}

			log.Infoln(log.Global, "Creating Engine ... Success.")

		} else {
			return engineInstance, errEngineCreated
		}
	} else {
		return engineInstance, errEngineCreated
	}
	return engineInstance, nil
}

// Start the engine. Starts the application and all services.
func (bot *Engine) Start() error {
	if bot == nil {
		return errNilEngine
	}
	log.Infoln(log.Global, "Starting Engine ...")

	// set current uptime to now
	engineLock.Lock()
	// ToDO: implement
	engineLock.Unlock()

	// start all services
	if err := bot.StartAll(); err != nil {
		return err
	}

	// Print some info
	cpus := bot.GoMaxProcessors
	if bot.GoMaxProcessors == -1 {
		cpus = runtime.NumCPU()
	}
	log.Infoln(log.Global, "Starting Engine ... Success.")
	log.Infof(log.Global, "Using %d out of %d logical processors.\n", cpus, runtime.NumCPU())
	return nil
}

// WaitForInterrupt is a blocking routine that returns when an operating system interrupt signal is received.
func (bot *Engine) WaitForInterrupt() error {
	if bot == nil {
		return errNilEngine
	}
	log.Infoln(log.Global, "Waiting for interrupt to shutdown ...")

	// main thread will block here until an interrupt is received
	interrupt := system.WaitForInterrupt()
	log.Infof(log.Global, "Captured '%v', requesting shutdown ...", interrupt)
	return nil
}

// Stop the engine. Stops all services and exits the application.
func (bot *Engine) Stop() error {
	if bot == nil {
		return errNilEngine
	}
	log.Infoln(log.Global, "Stopping Engine ...")

	// stop all services
	bot.StopAll()

	// wait for services to gracefully shutdown
	bot.ServiceManager.ServiceWG.Wait()

	log.Infoln(log.Global, "Stopping Engine ... Success.")

	// everything stopped, try closing log file
	if err := log.CloseLogger(); err != nil {
		return err
	}

	return nil
}

// Restart the engine. Try's to restart the application.
func (bot *Engine) Restart() error {
	if bot == nil {
		return errNilEngine
	}
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
