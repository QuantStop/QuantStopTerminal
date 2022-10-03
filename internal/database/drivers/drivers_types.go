package drivers

const (
	// DBSQLite const string for sqlite across code base
	DBSQLite = "sqlite"

	// DBSQLite3 const string for sqlite3 across code base
	DBSQLite3 = "sqlite3"

	// DBPostgreSQL const string for PostgreSQL across code base
	DBPostgreSQL = "postgres"

	// DBMySQL const string for MySQL across code base
	DBMySQL = "mysql"
)

// ConnectionDetails holds DSN information
type ConnectionDetails struct {
	Host     string
	Port     uint16
	Username string
	Password string
	Database string
	SSLMode  string
}
