package config

import (
	"encoding/json"
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/database"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/internal/system"
	"github.com/quantstop/quantstopterminal/internal/system/convert"
	jsonUtils "github.com/quantstop/quantstopterminal/internal/system/file/json"
	"github.com/quantstop/quantstopterminal/internal/webserver"
	"os"
	"path/filepath"
	"sync"
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
	LogConfig       *log.Config
	DatabaseConfig  *database.Config
	WebserverConfig *webserver.Config
}

func init() {
	findPaths()
}

// Refresh will rediscover the config paths
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

// NewConfig will create the Config object and set the default data paths for the application.
func NewConfig() (*Config, error) {

	config := &Config{}

	// A common use case is to get a private config folder for your app to
	// place its settings files into, that are specific to the local user.
	configPath := LocalConfig("QuantstopTerminal")
	err := makePath(configPath) // Ensure it exists.
	if err != nil {
		return nil, err
	}

	// Deal with a JSON configuration file in that folder.
	configFile := filepath.Join(configPath, "settings.json")

	// Does the file exist?
	if _, err = os.Stat(configFile); os.IsNotExist(err) {

		// Setup default config
		config.setupDefaultConfig(configPath)

		// Create the config file
		fh, err := os.Create(configFile)
		if err != nil {
			return nil, err
		}
		defer func(fh *os.File) {
			_ = fh.Close()
		}(fh)

		// Write config to file in json format
		err = jsonUtils.PrettyEncodeJson(&config, fh)
		if err != nil {
			return nil, err
		}

	} else {

		// Load the existing file.
		fh, err := os.Open(configFile)
		if err != nil {
			return nil, err
		}
		defer func(fh *os.File) {
			_ = fh.Close()
		}(fh)

		decoder := json.NewDecoder(fh)
		err = decoder.Decode(&config)
		if err != nil {
			return nil, err
		}
	}

	// Verify config
	if err = config.verifyConfig(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) setupDefaultConfig(configDir string) {

	c.ConfigDir = configDir
	c.GoMaxProcessors = -1

	// Generate default logging config
	c.LogConfig = log.NewConfig()

	// Copy default logging config to global log config
	log.RWM.Lock()
	log.GlobalLogConfig = c.LogConfig
	log.RWM.Unlock()

	// Generate default database config
	c.DatabaseConfig = database.NewConfig(configDir)

	// Generate default webserver config
	c.WebserverConfig = webserver.NewConfig(configDir)

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
		return err
	}

	// Verify config
	if err = c.verifyConfig(); err != nil {
		return fmt.Errorf("save config error: %v", err)
	}

	return nil
}

// verifyConfig will run Verify() on every service config to ensure all values are valid.
func (c *Config) verifyConfig() error {
	err := c.checkLoggerConfig()
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

	if c.LogConfig.Enabled == nil || c.LogConfig.Output == "" {
		c.LogConfig = log.NewConfig()
	}

	if c.LogConfig.AdvancedSettings.ShowLogSystemName == nil {
		c.LogConfig.AdvancedSettings.ShowLogSystemName = convert.BoolPtr(false)
	}

	if c.LogConfig.LoggerFileConfig != nil {
		if c.LogConfig.LoggerFileConfig.FileName == "" {
			c.LogConfig.LoggerFileConfig.FileName = "log.txt"
		}
		if c.LogConfig.LoggerFileConfig.Rotate == nil {
			c.LogConfig.LoggerFileConfig.Rotate = convert.BoolPtr(false)
		}
		if c.LogConfig.LoggerFileConfig.MaxSize <= 0 {
			//log.Warnf(log.Global, "Logger rotation size invalid, defaulting to %v", log.DefaultMaxFileSize)
			c.LogConfig.LoggerFileConfig.MaxSize = log.DefaultMaxFileSize
		}
		log.FileLoggingConfiguredCorrectly = true
	}
	log.RWM.Lock()
	log.GlobalLogConfig = c.LogConfig
	log.RWM.Unlock()

	logPath := c.GetDataPath("logs")
	err := system.CreateDir(logPath) //todo: what is this, it returns nil always?
	if err != nil {
		return err
	}
	log.FilePath = logPath

	return nil
}

// GetDataPath gets the data path for the given sub-path
func (c *Config) GetDataPath(elem ...string) string {
	return filepath.Join(append([]string{c.ConfigDir}, elem...)...)
}
