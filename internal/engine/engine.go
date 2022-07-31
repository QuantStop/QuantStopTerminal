package engine

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/quantstop/quantstopexchange/qsx"
	"github.com/quantstop/quantstopterminal/internal/config"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/internal/webserver"
	"github.com/quantstop/quantstopterminal/pkg/system"
	"github.com/quantstop/quantstopterminal/pkg/system/convert"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Engine struct {
	*Version
	Config              *config.Config
	SubsystemRegistry   *SubsystemRegistry
	DatabaseSubsystem   *DatabaseSubsystem
	NTPCheckerSubsystem *NTPCheckerSubsystem
	TraderSubsystem     *TraderSubsystem
	InternetSubsystem   *ConnectionMonitor
	SentimentAnalyzer   *SentimentAnalyzerSubsystem
	Webserver           *webserver.Webserver
	ExchangeManager     *ExchangeManager
	SubsystemWG         sync.WaitGroup
	Uptime              time.Time
}

const (
	DatabaseSubsystemName string = "database"
	NTPSubsystemName      string = "ntp_timekeeper"
	TraderSubsystemName   string = "active_trader"
	InternetCheckerName   string = "internet_monitor"
	SentimentAnalyzerName string = "sentiment_analyzer"
	ExchangeManagerName   string = "exchange_manager"
)

// engineMutex only locks and unlocks on engine creation functions
// as engine modifies global files, this protects the main bot creation
// functions from interfering with each other
var engineMutex sync.Mutex

// Create creates a new instance of the engine
func Create(config *config.Config, version *Version) (*Engine, error) {

	engineMutex.Lock()
	defer engineMutex.Unlock()

	if config == nil {
		return nil, errors.New("engine: config is nil")
	}

	var bot Engine
	var err error

	// Set the bot config
	bot.Config = config

	// Set the bot version
	bot.Version = version

	// Set the max processors for go
	err = system.AdjustGoMaxProcs(bot.Config.GoMaxProcessors)
	if err != nil {
		return nil, fmt.Errorf("unable to adjust runtime GOMAXPROCS value. Err: %s", err)
	}

	return &bot, nil
}

// Initialize sets up the engine, creating the subsystems, and the subsystem registry.
func (bot *Engine) Initialize() error {

	if bot == nil {
		return errors.New("engine instance is nil")
	}

	engineMutex.Lock()
	defer engineMutex.Unlock()

	// Create new subsystem registry
	bot.SubsystemRegistry = NewSubsystemRegistry()

	// Initialize database subsystem
	if err := bot.initDatabaseSubsystem(); err != nil {
		return err
	}

	// Initialize ntp checker subsystem
	if err := bot.initNtpMonitorSubsystem(); err != nil {
		return err
	}

	// Initialize strategy subsystem
	if err := bot.initStrategySubsystem(); err != nil {
		return err
	}

	// Initialize internet checker subsystem
	if err := bot.initInternetMonitorSubsystem(); err != nil {
		return err
	}

	// Initialize sentiment analyzer subsystem
	if err := bot.initSentimentAnalyzerSubsystem(); err != nil {
		return err
	}

	// Initialize exchange manager subsystem
	if err := bot.initExchangeManagerSubsystem(); err != nil {
		return err
	}

	return nil
}

func (bot *Engine) initDatabaseSubsystem() error {

	// Create and init database subsystem
	bot.DatabaseSubsystem = &DatabaseSubsystem{Subsystem: Subsystem{}}
	if err := bot.DatabaseSubsystem.init(bot, DatabaseSubsystemName); err != nil {
		log.Errorf(log.Global, "database subsystem unable to initialize: %v", err)
		return err
	}

	// Register database subsystem
	if err := bot.SubsystemRegistry.RegisterSubsystem(bot.DatabaseSubsystem); err != nil {
		log.Errorf(log.Global, "database subsystem unable to register: %v", err)
		return err
	}

	return nil
}

func (bot *Engine) initNtpMonitorSubsystem() error {
	if bot.Config.NTP.Enabled {

		// Create and init ntp checker subsystem
		bot.NTPCheckerSubsystem = &NTPCheckerSubsystem{Subsystem: Subsystem{}}
		if err := bot.NTPCheckerSubsystem.init(bot, NTPSubsystemName); err != nil {
			log.Errorf(log.Global, "NTP subsystem unable to initialize: %v", err)
			return err
		}

		// Register ntp checker subsystem
		if err := bot.SubsystemRegistry.RegisterSubsystem(bot.NTPCheckerSubsystem); err != nil {
			log.Errorf(log.Global, "NTP subsystem unable to register: %v", err)
			return err
		}

	}
	return nil
}

func (bot *Engine) initStrategySubsystem() error {
	/*if bot.Config.Strategy.Enabled {*/

	// Create and init strategy subsystem
	bot.TraderSubsystem = &TraderSubsystem{Subsystem: Subsystem{}}
	if err := bot.TraderSubsystem.init(bot, TraderSubsystemName); err != nil {
		log.Errorf(log.Global, "Trader subsystem unable to initialize: %v", err)
		return err
	}

	// Register strategy subsystem
	if err := bot.SubsystemRegistry.RegisterSubsystem(bot.TraderSubsystem); err != nil {
		log.Errorf(log.Global, "Trader subsystem unable to register: %v", err)
		return err
	}

	//}
	return nil
}

func (bot *Engine) initInternetMonitorSubsystem() error {
	if bot.Config.Internet.Enabled {

		// Create and init internet checker subsystem
		bot.InternetSubsystem = &ConnectionMonitor{Subsystem: Subsystem{}}
		if err := bot.InternetSubsystem.init(bot, InternetCheckerName); err != nil {
			log.Errorf(log.Global, "Internet checker subsystem unable to initialize: %v", err)
			return err
		}

		// Register internet checker subsystem
		if err := bot.SubsystemRegistry.RegisterSubsystem(bot.InternetSubsystem); err != nil {
			log.Errorf(log.Global, "Internet checker subsystem unable to register: %v", err)
			return err
		}

	}
	return nil
}

func (bot *Engine) initSentimentAnalyzerSubsystem() error {
	/*if bot.Config.Strategy.Enabled {*/

	// Create and init strategy subsystem
	bot.SentimentAnalyzer = &SentimentAnalyzerSubsystem{Subsystem: Subsystem{}}
	if err := bot.SentimentAnalyzer.init(bot, SentimentAnalyzerName); err != nil {
		log.Errorf(log.Global, "Sentiment Analyzer subsystem unable to initialize: %v", err)
		return err
	}

	// Register strategy subsystem
	if err := bot.SubsystemRegistry.RegisterSubsystem(bot.SentimentAnalyzer); err != nil {
		log.Errorf(log.Global, "Sentiment Analyzer subsystem unable to register: %v", err)
		return err
	}

	//}
	return nil
}

func (bot *Engine) initExchangeManagerSubsystem() error {

	// Create and init exchange manager subsystem
	bot.ExchangeManager = &ExchangeManager{Subsystem: Subsystem{}}
	if err := bot.ExchangeManager.init(bot, ExchangeManagerName); err != nil {
		log.Errorf(log.Global, "Exchange Manager subsystem unable to initialize: %v", err)
		return err
	}

	// Register exchange manager subsystem
	if err := bot.SubsystemRegistry.RegisterSubsystem(bot.ExchangeManager); err != nil {
		log.Errorf(log.Global, "Exchange Manager subsystem unable to register: %v", err)
		return err
	}

	return nil
}

// Run starts the newly created instance of the engine
func (bot *Engine) Run() error {

	if bot == nil {
		return errors.New("engine instance is nil")
	}

	engineMutex.Lock()
	defer engineMutex.Unlock()

	// Set the current uptime to now
	bot.Uptime = time.Now()

	// Start all subsystems in order of registration
	bot.SubsystemRegistry.StartAll(&bot.SubsystemWG)

	// Everything good, create and run webserver
	var err error
	bot.Webserver, err = webserver.CreateWebserver(bot, bot.Config.Webserver, bot.Version.IsDevelopment)
	if err != nil {
		return err
	}

	// Run the webserver (starts up both the http and websocket servers)
	go func() {
		err = bot.Webserver.ListenAndServe(true, bot.Config.ConfigDir)
		if err != nil {
			err = fmt.Errorf("unexpected error from ListenAndServe: %w", err)
			log.Error(log.Global, err)
		}
	}()

	// Run the trading subsystem
	err = bot.TraderSubsystem.run()
	if err != nil {
		return err
	}

	// Print some info
	log.Infof(log.Global, "QuantstopTerminal started.\n")
	log.Infof(log.Global,
		"Using %d out of %d logical processors for runtime performance\n",
		runtime.GOMAXPROCS(-1), runtime.NumCPU())

	return nil
}

// Stop stops the running instance of the engine
func (bot *Engine) Stop() {

	engineMutex.Lock()
	defer engineMutex.Unlock()

	log.Debugln(log.Global, "Engine shutting down..")

	// Stop webserver
	bot.Webserver.Shutdown()

	// Stop all subsystems
	bot.SubsystemRegistry.StopAll()

	// Wait for subsystems to gracefully shutdown
	bot.SubsystemWG.Wait()
	if err := log.CloseLogger(); err != nil {
		fmt.Printf("Failed to close logger. Error: %v\n", err)
	}

}

// Restart the running instance of the engine
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

// GetUptime returns the time since the bot last started
func (bot *Engine) GetUptime() string {
	//return time.Since(bot.Uptime).String()
	return convert.RoundDuration(time.Since(bot.Uptime), 2).String()
}

// SetConfig saves system configuration data
func (bot *Engine) SetConfig(apiUrl string, maxProcs string) error {
	intVar, err := strconv.Atoi(maxProcs)
	if err != nil {
		return err
	}
	bot.Config.GoMaxProcessors = intVar //todo: does this only take effect on restart?

	err = bot.Config.SaveConfig()
	if err != nil {
		return err
	}

	err = bot.Restart()
	if err != nil {
		return err
	}

	return nil
}

// GetSubsystemsStatus returns the status of all engine subsystems
func (bot *Engine) GetSubsystemsStatus() map[string]bool {

	status := make(map[string]bool)

	if bot.DatabaseSubsystem == nil {
		status[DatabaseSubsystemName] = false
	} else {
		status[DatabaseSubsystemName] = bot.DatabaseSubsystem.isRunning()
	}

	if bot.NTPCheckerSubsystem == nil {
		status[NTPSubsystemName] = false
	} else {
		status[NTPSubsystemName] = bot.NTPCheckerSubsystem.isRunning()
	}

	if bot.TraderSubsystem == nil {
		status[TraderSubsystemName] = false
	} else {
		status[TraderSubsystemName] = bot.TraderSubsystem.isRunning()
	}

	if bot.InternetSubsystem == nil {
		status[InternetCheckerName] = false
	} else {
		status[InternetCheckerName] = bot.InternetSubsystem.isRunning()
	}

	if bot.DatabaseSubsystem == nil {
		status[DatabaseSubsystemName] = false
	} else {
		status[DatabaseSubsystemName] = bot.DatabaseSubsystem.isRunning()
	}

	return status
}

// SetSubsystem enables or disables an engine subsystem
func (bot *Engine) SetSubsystem(subSystemName string, enable bool) error {
	if bot == nil {
		return errors.New("engine instance is nil")
	}

	if bot.Config == nil {
		return errNilEngine
	}

	engineMutex.Lock()
	defer engineMutex.Unlock()

	var err error
	switch strings.ToLower(subSystemName) {

	case NTPSubsystemName:
		if enable {
			if bot.NTPCheckerSubsystem == nil {
				err = bot.NTPCheckerSubsystem.init(bot, NTPSubsystemName)
				if err != nil {
					return err
				}
			}
			return bot.NTPCheckerSubsystem.start(&bot.SubsystemWG)
		} else {
			return bot.NTPCheckerSubsystem.stop()
		}

	case TraderSubsystemName:
		if enable {
			if bot.TraderSubsystem == nil {
				err = bot.TraderSubsystem.init(bot, TraderSubsystemName)
				if err != nil {
					return err
				}
			}
			return bot.TraderSubsystem.start(&bot.SubsystemWG)
		} else {
			return bot.TraderSubsystem.stop()
		}

	case InternetCheckerName:
		if enable {
			if bot.InternetSubsystem == nil {
				err = bot.InternetSubsystem.init(bot, InternetCheckerName)
				if err != nil {
					return err
				}
			}
			return bot.InternetSubsystem.start(&bot.SubsystemWG)
		} else {
			return bot.InternetSubsystem.stop()
		}

	}
	return fmt.Errorf("%s: %w", subSystemName, ErrSubsystemNotFound)
}

// GetVersion returns a map of the current version, along with other info
func (bot *Engine) GetVersion() map[string]string {
	version := make(map[string]string)

	version["version"] = bot.Version.Version
	version["copyright"] = bot.Version.Copyright
	version["prereleaseblurb"] = bot.Version.PrereleaseBlurb
	version["github"] = bot.Version.GitHub
	version["issues"] = bot.Version.Issues
	if bot.Version.IsDaemon {
		version["isdaemon"] = "true"
	} else {
		version["isdaemon"] = "false"
	}
	if bot.Version.IsRelease {
		version["isrelease"] = "true"
	} else {
		version["isrelease"] = "false"
	}
	if bot.Version.IsDevelopment {
		version["isdevelopment"] = "true"
	} else {
		version["isdevelopment"] = "false"
	}

	return version

}

// GetSQL returns a pointer to the database connection for the given database name
func (bot *Engine) GetSQL(dbName string) (*sql.DB, error) {
	switch dbName {
	case "core":
		if bot.DatabaseSubsystem.coreDatabase.SQL != nil {
			return bot.DatabaseSubsystem.coreDatabase.SQL, nil
		}
		log.Errorln(log.Global, ErrNilCoreSQL)
		return nil, ErrNilCoreSQL
	case "tda":
		if bot.DatabaseSubsystem.tdameritradeDatabase.SQL != nil {
			return bot.DatabaseSubsystem.tdameritradeDatabase.SQL, nil
		}
		log.Errorln(log.Global, ErrNilTDAmeritradeSQL)
		return nil, ErrNilTDAmeritradeSQL
	case "coinbase":
		if bot.DatabaseSubsystem.coinbaseDatabase.SQL != nil {
			return bot.DatabaseSubsystem.coinbaseDatabase.SQL, nil
		}
		log.Errorln(log.Global, ErrNilCoinbaseSQL)
		return nil, ErrNilCoinbaseSQL
	default:
		return nil, errors.New("GetSQL unknown database name provided")
	}
}

// GetExchange returns an exchange connection
func (bot *Engine) GetExchange(name string) qsx.IExchange {
	switch name {
	case "coinbasepro":
		return bot.ExchangeManager.Exchanges["coinbasepro"]
	}
	return nil
}

// GetSupportedExchangesList returns a list of all supported exchanges
func (bot *Engine) GetSupportedExchangesList() []string {
	var list []string
	for _, e := range bot.ExchangeManager.Exchanges {
		list = append(list, string(e.GetName()))
	}
	return list
}
