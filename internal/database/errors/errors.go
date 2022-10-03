package errors

import "errors"

var (
	// ErrNoDatabaseProvided error to display when no database is provided
	ErrNoDatabaseProvided = errors.New("no database provided")

	// ErrDatabaseSupportDisabled error to display when no database is provided
	ErrDatabaseSupportDisabled = errors.New("database support is disabled")

	// ErrFailedToConnect for when a database fails to connect
	ErrFailedToConnect = errors.New("database failed to connect")

	// ErrDatabaseNotConnected for when a database is not connected
	ErrDatabaseNotConnected = errors.New("database is not connected")

	// ErrNilInstance for when a database is nil
	ErrNilInstance = errors.New("database instance is nil")

	// ErrNilConfig for when a config is nil
	ErrNilConfig = errors.New("received nil config")

	ErrNilSQL = errors.New("database SQL connection is nil")

	ErrFailedPing = errors.New("unable to verify database is connected, failed ping")
)
