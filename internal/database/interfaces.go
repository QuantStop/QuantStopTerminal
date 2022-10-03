package database

import (
	"context"
	"database/sql"
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
