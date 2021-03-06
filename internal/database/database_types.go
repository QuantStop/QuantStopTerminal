package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/quantstop/quantstopterminal/internal/database/drivers"
	"sync"
)

// Instance holds all information for a database instance
type Instance struct {
	SQL       *sql.DB
	DataPath  string
	config    *Config
	connected bool
	m         sync.RWMutex
}

// Config holds all configurable parameters for the database subsystem
type Config struct {
	Enabled bool
	Verbose bool
	Driver  string
	DSN     drivers.ConnectionDetails
}

var (
	// CoreDB Global CoreDB Connection
	CoreDB = &Instance{}

	// CoinbaseDB Global CoreDB Connection
	CoinbaseDB = &Instance{}

	// TDAmeritradeDB Global CoreDB Connection
	TDAmeritradeDB = &Instance{}

	// MigrationDir which folder to look in for current migrations
	//MigrationDir = filepath.Join("..", "..", "database", "migrations")

	// ErrNoDatabaseProvided error to display when no database is provided
	ErrNoDatabaseProvided = errors.New("no database provided")

	// ErrDatabaseSupportDisabled error to display when no database is provided
	ErrDatabaseSupportDisabled = errors.New("database support is disabled")

	// SupportedDrivers slice of supported database driver types
	SupportedDrivers = []string{DBSQLite, DBSQLite3, DBPostgreSQL, DBMySQL}

	// ErrFailedToConnect for when a database fails to connect
	ErrFailedToConnect = errors.New("database failed to connect")

	// ErrDatabaseNotConnected for when a database is not connected
	ErrDatabaseNotConnected = errors.New("database is not connected")

	// ErrNilInstance for when a database is nil
	ErrNilInstance = errors.New("database instance is nil")

	// ErrNilConfig for when a config is nil
	ErrNilConfig  = errors.New("received nil config")
	errNilSQL     = errors.New("database SQL connection is nil")
	errFailedPing = errors.New("unable to verify database is connected, failed ping")
)

const (
	// DBSQLite const string for sqlite across code base
	DBSQLite = "sqlite"

	// DBSQLite3 const string for sqlite3 across code base
	DBSQLite3 = "sqlite3"

	// DBPostgreSQL const string for PostgreSQL across code base
	DBPostgreSQL = "postgres"

	// DBMySQL const string for MySQL across code base
	DBMySQL = "mysql"

	// DefaultCoreDatabase const string for name of core database (sqlite filename)
	DefaultCoreDatabase = "qst.db"

	// DefaultCoinbaseDatabase const string for name of coinbase database (sqlite filename)
	DefaultCoinbaseDatabase = "qst_coinbase.db"

	// DefaultTDAmeritradeDatabase const string for name of td-ameritrade database (sqlite filename)
	DefaultTDAmeritradeDatabase = "qst_tdameritrade.db"
)

// IDatabase allows for the passing of a database struct
// without giving the receiver access to all functionality
type IDatabase interface {
	IsConnected() bool
	GetSQL() (*sql.DB, error)
	GetConfig() *Config
}

// ISQL allows for the passing of an SQL connection
// without giving the receiver access to all functionality
type ISQL interface {
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func GenDefaultSettings(name string) *Config {

	switch name {
	case "core":
		return &Config{
			Enabled: true,
			Verbose: true,
			Driver:  "sqlite",
			DSN: drivers.ConnectionDetails{
				Host:     "127.0.0.1",
				Port:     3306,
				Username: "docker",
				Password: "docker",
				Database: DefaultCoreDatabase,
				SSLMode:  "false",
			},
		}
	case "coinbase":
		return &Config{
			Enabled: true,
			Verbose: true,
			Driver:  "sqlite",
			DSN: drivers.ConnectionDetails{
				Host:     "127.0.0.1",
				Port:     3306,
				Username: "docker",
				Password: "docker",
				Database: DefaultCoinbaseDatabase,
				SSLMode:  "false",
			},
		}
	case "tdameritrade":
		return &Config{
			Enabled: true,
			Verbose: true,
			Driver:  "sqlite",
			DSN: drivers.ConnectionDetails{
				Host:     "127.0.0.1",
				Port:     3306,
				Username: "docker",
				Password: "docker",
				Database: DefaultTDAmeritradeDatabase,
				SSLMode:  "false",
			},
		}
	default:
		return &Config{
			Enabled: true,
			Verbose: true,
			Driver:  "sqlite",
			DSN: drivers.ConnectionDetails{
				Host:     "127.0.0.1",
				Port:     3306,
				Username: "docker",
				Password: "docker",
				Database: DefaultCoreDatabase,
				SSLMode:  "false",
			},
		}
	}

}
