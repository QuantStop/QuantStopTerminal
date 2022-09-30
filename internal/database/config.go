package database

const (
	// DBSQLite const string for sqlite across code base
	DBSQLite = "sqlite"

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
	CoreConn *ConnectionDetails
}

// ConnectionDetails holds DSN information for a single database connection
type ConnectionDetails struct {
	Driver   string
	Host     string
	Port     uint16
	Username string
	Password string
	Database string
	SSLMode  string
}

// NewConfig Generate default settings for the Config struct
func NewConfig() *Config {
	return &Config{
		CoreConn: &ConnectionDetails{
			Driver:   DBSQLite,
			Host:     DefaultHost,
			Port:     DefaultPort,
			Username: DefaultUsername,
			Password: DefaultPassword,
			Database: DefaultCoreDatabase,
			SSLMode:  DefaultSSLMode,
		},
	}
}
