package database

import (
	"github.com/quantstop/quantstopterminal/internal/database/drivers"
	"github.com/quantstop/quantstopterminal/internal/database/repository"
)

const (

	// DefaultCoreDatabase const string for name of core database (sqlite filename)
	DefaultCoreDatabase = "qst.db"

	// DefaultCoinbaseDatabase const string for name of coinbase database (sqlite filename)
	DefaultCoinbaseDatabase = "qst_coinbase.db"

	// DefaultTDAmeritradeDatabase const string for name of td-ameritrade database (sqlite filename)
	DefaultTDAmeritradeDatabase = "qst_tdameritrade.db"

	// DefaultHost the default host to use to connect to the database
	DefaultHost = "127.0.0.1"

	// DefaultPort the default port to use to connect to the database
	DefaultPort = 3306

	// DefaultUsername the default username to use when first creating the database
	DefaultUsername = "qst"

	// DefaultPassword the default password to use when first creating the database
	DefaultPassword = "qu4n75t0pt3rm1nal1s4w3s0m3!"

	// DefaultSSLMode default ssl mode is off
	DefaultSSLMode = "false"
)

// Config holds information for all databases
type Config struct {
	CoreConfig         *repository.InstanceConfig
	CoinbaseConfig     *repository.InstanceConfig
	TDAmeritradeConfig *repository.InstanceConfig
}

// NewConfig Generate default settings for the Config struct
func NewConfig(configDir string) *Config {
	return &Config{
		CoreConfig: &repository.InstanceConfig{
			Enabled:   true,
			Verbose:   true,
			Driver:    drivers.DBSQLite,
			ConfigDir: configDir,
			DSN: &drivers.ConnectionDetails{
				Host:     DefaultHost,
				Port:     DefaultPort,
				Username: DefaultUsername,
				Password: DefaultPassword,
				Database: DefaultCoreDatabase,
				SSLMode:  DefaultSSLMode,
			},
		},
		CoinbaseConfig: &repository.InstanceConfig{
			Enabled:   true,
			Verbose:   true,
			Driver:    drivers.DBSQLite,
			ConfigDir: configDir,
			DSN: &drivers.ConnectionDetails{
				Host:     DefaultHost,
				Port:     DefaultPort,
				Username: DefaultUsername,
				Password: DefaultPassword,
				Database: DefaultCoinbaseDatabase,
				SSLMode:  DefaultSSLMode,
			},
		},
		TDAmeritradeConfig: &repository.InstanceConfig{
			Enabled:   true,
			Verbose:   true,
			Driver:    drivers.DBSQLite,
			ConfigDir: configDir,
			DSN: &drivers.ConnectionDetails{
				Host:     DefaultHost,
				Port:     DefaultPort,
				Username: DefaultUsername,
				Password: DefaultPassword,
				Database: DefaultTDAmeritradeDatabase,
				SSLMode:  DefaultSSLMode,
			},
		},
	}
}
