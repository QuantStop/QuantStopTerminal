package config

import (
	"encoding/json"
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/connectionmonitor"
	"github.com/quantstop/quantstopterminal/internal/database"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/internal/ntpmonitor"
	"github.com/quantstop/quantstopterminal/internal/webserver"
	"github.com/quantstop/quantstopterminal/pkg/system"
	"github.com/quantstop/quantstopterminal/pkg/system/convert"
	jsonUtils "github.com/quantstop/quantstopterminal/pkg/system/file/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	DefaultNTPAllowedDifference         = 50000000
	DefaultNTPAllowedNegativeDifference = 50000000
	DefaultFileMode                     = os.FileMode(0755)
)

var (
	mutex sync.Mutex
)

type Config struct {
	ConfigDir       string
	GoMaxProcessors int
	FirstLoginSetup bool
	CoreDB          database.Config
	CoinbaseDB      database.Config
	TDAmeritradeDB  database.Config
	Webserver       *webserver.Config
	NTP             ntpmonitor.Config
	Internet        connectionmonitor.Config
	Logger          log.Config
}

func init() {
	findPaths()
}

// Refresh will rediscover the config paths, checking current environment
// variables again.
//
// This function is automatically called when the program initializes. If you
// change the environment variables at run-time, though, you may call the
// Refresh() function to reevaluate the config paths.
func Refresh() {
	findPaths()
}

// SystemConfig returns the system-wide configuration paths, with optional path
// components added to the end for vendor/application-specific settings.
func SystemConfig(folder ...string) []string {
	if len(folder) == 0 {
		return systemConfig
	}

	var paths []string
	for _, root := range systemConfig {
		p := append([]string{root}, filepath.Join(folder...))
		paths = append(paths, filepath.Join(p...))
	}

	return paths
}

// LocalConfig returns the local user configuration path, with optional
// path components added to the end for vendor/application-specific settings.
func LocalConfig(folder ...string) string {
	if len(folder) == 0 {
		return localConfig
	}

	return filepath.Join(localConfig, filepath.Join(folder...))
}

// LocalCache returns the local user cache folder, with optional path
// components added to the end for vendor/application-specific settings.
func LocalCache(folder ...string) string {
	if len(folder) == 0 {
		return localCache
	}

	return filepath.Join(localCache, filepath.Join(folder...))
}

// makePath ensures that the full path you wanted, including vendor or
// application-specific components, exists. You can give this the output of
// any config path functions (SystemConfig, LocalConfig or LocalCache).
//
// In the event that the path function gives multiple answers, e.g. for
// SystemConfig, MakePath() will only attempt to create the sub-folders on
// the *first* path found. If this isn't what you want, you may want to just
// use the os.MkdirAll() functionality directly.
func makePath(paths ...string) error {
	if len(paths) >= 1 {
		err := os.MkdirAll(paths[0], DefaultFileMode)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetupConfig will create the Config object and set the default data paths for the application.
func (c *Config) SetupConfig() error {

	// A common use case is to get a private config folder for your app to
	// place its settings files into, that are specific to the local user.
	configPath := LocalConfig("QuantstopTerminal")
	err := makePath(configPath) // Ensure it exists.
	if err != nil {
		return err
	}

	// Deal with a JSON configuration file in that folder.
	configFile := filepath.Join(configPath, "settings.json")

	// Does the file not exist?
	if _, err = os.Stat(configFile); os.IsNotExist(err) {

		// Setup default config
		c.ConfigDir = configPath
		c.GoMaxProcessors = -1
		c.FirstLoginSetup = true
		c.CoreDB = *database.GenDefaultSettings("core")
		c.CoinbaseDB = *database.GenDefaultSettings("coinbase")
		c.TDAmeritradeDB = *database.GenDefaultSettings("tdameritrade")
		c.Webserver = &webserver.Config{
			Enabled:             true,
			HttpListenAddr:      ":443",
			WebsocketListenAddr: ":8090",
		}
		c.NTP = ntpmonitor.Config{
			Enabled: true,
			Verbose: false,
			Level:   0,
			Pool: []string{
				"pool.ntp.org:123",
			},
			AllowedDifference:         new(time.Duration),
			AllowedNegativeDifference: new(time.Duration),
		}
		c.Internet = connectionmonitor.Config{
			Enabled:          true,
			Initialized:      false,
			DNSList:          []string{"8.8.8.8", "8.8.4.4", "1.1.1.1", "1.0.0.1"},
			PublicDomainList: []string{"www.google.com", "www.cloudflare.com", "www.facebook.com"},
			CheckInterval:    time.Second * 3,
		}

		// Set default ntp settings
		*c.NTP.AllowedDifference = DefaultNTPAllowedDifference
		*c.NTP.AllowedNegativeDifference = DefaultNTPAllowedNegativeDifference

		// Load default logging config
		c.Logger = *log.GenDefaultSettings()

		// Copy default logging config to global log config
		log.RWM.Lock()
		log.GlobalLogConfig = &c.Logger
		log.RWM.Unlock()

		// Create the config file
		fh, err := os.Create(configFile)
		if err != nil {
			return err
		}
		defer func(fh *os.File) {
			_ = fh.Close()
		}(fh)

		// Write config to file in json format
		err = jsonUtils.PrettyEncodeJson(&c, fh)
		if err != nil {
			//log.Fatal(err)
			log.Error(log.Global, err)
		}

	} else {
		// Load the existing file.
		fh, err := os.Open(configFile)
		if err != nil {
			return err
		}
		defer func(fh *os.File) {
			_ = fh.Close()
		}(fh)

		decoder := json.NewDecoder(fh)
		err = decoder.Decode(&c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) SaveConfig() error {
	// A common use case is to get a private config folder for your app to
	// place its settings files into, that are specific to the local user.
	configPath := LocalConfig("QuantstopTerminal")

	// Deal with a JSON configuration file in that folder.
	configFile := filepath.Join(configPath, "settings.json")
	//fh, err := os.Open(configFile)
	fh, err := os.OpenFile(configFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, DefaultFileMode)
	if err != nil {
		return err
	}
	defer func(fh *os.File) {
		_ = fh.Close()
	}(fh)
	// Write config to file in json format
	err = jsonUtils.PrettyEncodeJson(&c, fh)
	if err != nil {
		//log.Fatal(err)
		log.Error(log.Global, err)
		return err
	}
	return nil
}

// CheckConfig will run private functions to verify the system config, and all subsystem configs are valid
func (c *Config) CheckConfig() error {
	err := c.checkLoggerConfig()
	if err != nil {
		return err
	}

	err = c.checkDatabaseConfig()
	if err != nil {
		return err
	}

	return nil
}

// CheckLoggerConfig checks to see logger values are present and valid in config
// if not, it creates a default instance of the logger
func (c *Config) checkLoggerConfig() error {
	mutex.Lock()
	defer mutex.Unlock()

	if c.Logger.Enabled == nil || c.Logger.Output == "" {
		c.Logger = *log.GenDefaultSettings()
	}

	if c.Logger.AdvancedSettings.ShowLogSystemName == nil {
		c.Logger.AdvancedSettings.ShowLogSystemName = convert.BoolPtr(false)
	}

	if c.Logger.LoggerFileConfig != nil {
		if c.Logger.LoggerFileConfig.FileName == "" {
			c.Logger.LoggerFileConfig.FileName = "log.txt"
		}
		if c.Logger.LoggerFileConfig.Rotate == nil {
			c.Logger.LoggerFileConfig.Rotate = convert.BoolPtr(false)
		}
		if c.Logger.LoggerFileConfig.MaxSize <= 0 {
			log.Warnf(log.Global, "Logger rotation size invalid, defaulting to %v", log.DefaultMaxFileSize)
			c.Logger.LoggerFileConfig.MaxSize = log.DefaultMaxFileSize
		}
		log.FileLoggingConfiguredCorrectly = true
	}
	log.RWM.Lock()
	log.GlobalLogConfig = &c.Logger
	log.RWM.Unlock()

	logPath := c.GetDataPath("logs")
	err := system.CreateDir(logPath)
	if err != nil {
		return err
	}
	log.LogPath = logPath

	return nil
}

func (c *Config) checkDatabaseConfig() error {
	mutex.Lock()
	defer mutex.Unlock()

	// todo: make work for all databases

	if (c.CoreDB == database.Config{}) {
		c.CoreDB.Driver = database.DBSQLite3
		c.CoreDB.DSN.Database = database.DefaultCoreDatabase
	}

	if !c.CoreDB.Enabled {
		return nil
	}

	if !system.StringDataCompare(database.SupportedDrivers, c.CoreDB.Driver) {
		c.CoreDB.Enabled = false
		return fmt.Errorf("unsupported database driver %v, database disabled", c.CoreDB.Driver)
	}

	if c.CoreDB.Driver == database.DBSQLite || c.CoreDB.Driver == database.DBSQLite3 {
		databaseDir := c.GetDataPath("database")
		err := system.CreateDir(databaseDir)
		if err != nil {
			return err
		}
		database.CoreDB.DataPath = databaseDir
	}

	return database.CoreDB.SetConfig(&c.CoreDB)
}

// GetDataPath gets the data path for the given sub-path
func (c *Config) GetDataPath(elem ...string) string {
	return filepath.Join(append([]string{c.ConfigDir}, elem...)...)
}
