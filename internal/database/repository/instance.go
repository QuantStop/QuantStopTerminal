package repository

import (
	"database/sql"
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/database/drivers"
	"github.com/quantstop/quantstopterminal/internal/database/errors"
	"sync"
	"time"
)

var (
	// CoreDB Global CoreDB Connection
	CoreDB = &Instance{}

	// CoinbaseDB Global CoreDB Connection
	CoinbaseDB = &Instance{}

	// TDAmeritradeDB Global CoreDB Connection
	TDAmeritradeDB = &Instance{}
)

type InstanceConfig struct {
	Enabled   bool
	Verbose   bool
	Driver    string
	ConfigDir string
	DSN       *drivers.ConnectionDetails
}

// Instance holds all information for a single database instance
type Instance struct {
	SQL       *sql.DB
	config    *InstanceConfig
	connected bool
	m         sync.RWMutex
}

// SetConfig safely sets the global database instance's config
func (i *Instance) SetConfig(cfg *InstanceConfig) error {
	if i == nil {
		return errors.ErrNilInstance
	}
	if cfg == nil {
		return errors.ErrNilConfig
	}
	i.m.Lock()
	i.config = cfg
	i.m.Unlock()
	return nil
}

// SetSQLiteConnection safely sets the global database instance's connection to use SQLite
func (i *Instance) SetSQLiteConnection(con *sql.DB) error {
	if i == nil {
		return errors.ErrNilInstance
	}
	if con == nil {
		return errors.ErrNilSQL
	}
	i.m.Lock()
	defer i.m.Unlock()
	i.SQL = con
	i.SQL.SetMaxOpenConns(1)
	return nil
}

// SetPostgresConnection safely sets the global database instance's connection to use Postgres
func (i *Instance) SetPostgresConnection(con *sql.DB) error {
	if i == nil {
		return errors.ErrNilInstance
	}
	if con == nil {
		return errors.ErrNilSQL
	}
	if err := con.Ping(); err != nil {
		return fmt.Errorf("%w %s", errors.ErrFailedPing, err)
	}
	i.m.Lock()
	defer i.m.Unlock()
	i.SQL = con
	i.SQL.SetMaxOpenConns(2)
	i.SQL.SetMaxIdleConns(1)
	i.SQL.SetConnMaxLifetime(time.Hour)
	return nil
}

// SetMySQLConnection safely sets the global database instance's connection to use SQLite
func (i *Instance) SetMySQLConnection(con *sql.DB) error {
	if i == nil {
		return errors.ErrNilInstance
	}
	if con == nil {
		return errors.ErrNilSQL
	}
	i.m.Lock()
	defer i.m.Unlock()
	i.SQL = con
	i.SQL.SetMaxOpenConns(10)
	i.SQL.SetMaxIdleConns(10)
	i.SQL.SetConnMaxLifetime(time.Hour)
	return nil
}

// SetConnected safely sets the global database instance's connected status
func (i *Instance) SetConnected(v bool) {
	if i == nil {
		return
	}
	i.m.Lock()
	i.connected = v
	i.m.Unlock()
}

// CloseConnection safely disconnects the global database instance
func (i *Instance) CloseConnection() error {
	if i == nil {
		return errors.ErrNilInstance
	}
	if i.SQL == nil {
		return errors.ErrNilSQL
	}
	i.m.Lock()
	defer i.m.Unlock()

	return i.SQL.Close()
}

// IsConnected safely checks the SQL connection status
func (i *Instance) IsConnected() bool {
	if i == nil {
		return false
	}
	i.m.RLock()
	defer i.m.RUnlock()
	return i.connected
}

// GetConfig safely returns a copy of the config
func (i *Instance) GetConfig() *InstanceConfig {
	if i == nil {
		return nil
	}
	i.m.RLock()
	defer i.m.RUnlock()
	cpy := i.config
	return cpy
}

// Ping pings the database
func (i *Instance) Ping() error {
	if i == nil {
		return errors.ErrNilInstance
	}
	if !i.IsConnected() {
		return errors.ErrDatabaseNotConnected
	}
	i.m.RLock()
	defer i.m.RUnlock()
	if i.SQL == nil {
		return errors.ErrNilSQL
	}
	return i.SQL.Ping()
}

// GetSQL returns the sql connection
func (i *Instance) GetSQL() (*sql.DB, error) {
	if i == nil {
		return nil, errors.ErrNilInstance
	}
	if i.SQL == nil {
		return nil, errors.ErrNilSQL
	}
	i.m.Lock()
	defer i.m.Unlock()
	resp := i.SQL
	return resp, nil
}
