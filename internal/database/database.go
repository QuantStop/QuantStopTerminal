package database

import (
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/database/drivers"
	"github.com/quantstop/quantstopterminal/internal/database/drivers/mysql"
	"github.com/quantstop/quantstopterminal/internal/database/drivers/postgres"
	"github.com/quantstop/quantstopterminal/internal/database/drivers/sqlite3"
	"github.com/quantstop/quantstopterminal/internal/database/errors"
	"github.com/quantstop/quantstopterminal/internal/database/repository"
	"github.com/quantstop/quantstopterminal/internal/database/seed"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/internal/service"

	"sync"
	"time"
)

var (

// MigrationDir which folder to look in for current migrations
//MigrationDir = filepath.Join("..", "..", "database", "migrations")

// SupportedDrivers slice of supported database driver types
// SupportedDrivers = []string{DBSQLite, DBSQLite3, DBPostgreSQL, DBMySQL}
)

type Database struct {
	*Config
	*service.Service
	CoreDB         *repository.CoreDatabase
	CoinbaseDB     *repository.Instance
	TDAmeritradeDB *repository.Instance
}

// NewDatabase creates the database connections
func NewDatabase(config *Config) (*Database, error) {

	if config == nil {
		return nil, fmt.Errorf("create database failed, nil Config received")
	}

	db := &Database{
		Service: service.NewService("database", true),
		Config:  config,
	}

	var err error

	if config.CoreConfig == nil {
		return nil, fmt.Errorf("create database failed, nil CoreConfig received")
	}

	if config.CoreConfig.Enabled {
		db.CoreDB = &repository.CoreDatabase{Instance: repository.CoreDB}
		if err = db.CoreDB.SetConfig(config.CoreConfig); err != nil {
			return nil, err
		}
		switch config.CoreConfig.Driver {
		case drivers.DBPostgreSQL:
			log.Debugf(log.Database,
				"database subsystem attempting to establish database connection to host %s/%s utilising %s driver\n",
				config.CoreConfig.DSN.Host,
				config.CoreConfig.DSN.Database,
				config.CoreConfig.Driver)
			db.CoreDB.Instance, err = postgres.Connect("core", config.CoreConfig.DSN)
		case drivers.DBSQLite, drivers.DBSQLite3:
			log.Debugf(log.Database,
				"database subsystem attempting to establish database connection to %s utilising %s driver\n",
				config.CoreConfig.DSN.Database,
				config.CoreConfig.Driver)
			db.CoreDB.Instance, err = sqlite.Connect("core", config.CoreConfig.DSN.Database)
		case drivers.DBMySQL:
			log.Debugf(log.Database,
				"database subsystem attempting to establish database connection to host %s/%s utilising %s driver\n",
				config.CoreConfig.DSN.Host,
				config.CoreConfig.DSN.Database,
				config.CoreConfig.Driver)
			db.CoreDB.Instance, err = mysql.Connect("core", config.CoreConfig.DSN)
		default:
			return nil, errors.ErrNoDatabaseProvided
		}
		if err != nil {
			return nil, fmt.Errorf("%w: %v Core database unavailable", errors.ErrFailedToConnect, err)
		}
		db.CoreDB.SetConnected(true)
	}

	if config.CoinbaseConfig.Enabled {
		db.CoinbaseDB = repository.CoinbaseDB
		if err = db.CoinbaseDB.SetConfig(config.CoinbaseConfig); err != nil {
			return nil, err
		}
		switch config.CoinbaseConfig.Driver {
		case drivers.DBPostgreSQL:
			log.Debugf(log.Database,
				"database subsystem attempting to establish database connection to host %s/%s utilising %s driver\n",
				config.CoinbaseConfig.DSN.Host,
				config.CoinbaseConfig.DSN.Database,
				config.CoinbaseConfig.Driver)
			db.CoinbaseDB, err = postgres.Connect("coinbase", config.CoinbaseConfig.DSN)
		case drivers.DBSQLite, drivers.DBSQLite3:
			log.Debugf(log.Database,
				"database subsystem attempting to establish database connection to %s utilising %s driver\n",
				config.CoinbaseConfig.DSN.Database,
				config.CoinbaseConfig.Driver)
			db.CoinbaseDB, err = sqlite.Connect("coinbase", config.CoinbaseConfig.DSN.Database)
		case drivers.DBMySQL:
			log.Debugf(log.Database,
				"database subsystem attempting to establish database connection to host %s/%s utilising %s driver\n",
				config.CoinbaseConfig.DSN.Host,
				config.CoinbaseConfig.DSN.Database,
				config.CoinbaseConfig.Driver)
			db.CoinbaseDB, err = mysql.Connect("coinbase", config.CoinbaseConfig.DSN)
		default:
			return nil, errors.ErrNoDatabaseProvided
		}
		if err != nil {
			return nil, fmt.Errorf("%w: %v Coinbase database unavailable", errors.ErrFailedToConnect, err)
		}
		db.CoinbaseDB.SetConnected(true)
	}

	if config.TDAmeritradeConfig.Enabled {
		db.TDAmeritradeDB = repository.TDAmeritradeDB
		if err = db.TDAmeritradeDB.SetConfig(config.TDAmeritradeConfig); err != nil {
			return nil, err
		}
		switch config.TDAmeritradeConfig.Driver {
		case drivers.DBPostgreSQL:
			log.Debugf(log.Database,
				"database subsystem attempting to establish database connection to host %s/%s utilising %s driver\n",
				config.TDAmeritradeConfig.DSN.Host,
				config.TDAmeritradeConfig.DSN.Database,
				config.TDAmeritradeConfig.Driver)
			db.TDAmeritradeDB, err = postgres.Connect("tdameritrade", config.TDAmeritradeConfig.DSN)
		case drivers.DBSQLite, drivers.DBSQLite3:
			log.Debugf(log.Database,
				"database subsystem attempting to establish database connection to %s utilising %s driver\n",
				config.TDAmeritradeConfig.DSN.Database,
				config.TDAmeritradeConfig.Driver)
			db.TDAmeritradeDB, err = sqlite.Connect("tdameritrade", config.TDAmeritradeConfig.DSN.Database)
		case drivers.DBMySQL:
			log.Debugf(log.Database,
				"database subsystem attempting to establish database connection to host %s/%s utilising %s driver\n",
				config.TDAmeritradeConfig.DSN.Host,
				config.TDAmeritradeConfig.DSN.Database,
				config.TDAmeritradeConfig.Driver)
			db.TDAmeritradeDB, err = mysql.Connect("tdameritrade", config.TDAmeritradeConfig.DSN)
		default:
			return nil, errors.ErrNoDatabaseProvided
		}
		if err != nil {
			return nil, fmt.Errorf("%w: %v TD-Ameritrade database unavailable", errors.ErrFailedToConnect, err)
		}
		db.TDAmeritradeDB.SetConnected(true)
	}

	// todo: seed all databases
	err = seed.DatabaseSeed(db.CoreDB.SQL, db.CoreConfig.Driver)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Start spawns all processes done by the service.
func (db *Database) Start(group *sync.WaitGroup) error {
	if err := db.Service.Start(group); err != nil {
		return err
	}
	go db.Run(group)
	log.Debugln(log.Webserver, db.GetName()+service.MsgServiceStarted)
	return nil
}

// Run is the main thread of the process, called as a goroutine.
func (db *Database) Run(wg *sync.WaitGroup) {
	t := time.NewTicker(time.Second * 5)

	// This function runs when the for loop returns
	defer func() {
		t.Stop()
		wg.Done()
		log.Debugln(log.Database, db.GetName()+service.MsgServiceShutdown)
	}()

	// This lets the goroutine wait on communication from the channel
	// Docs: https://tour.golang.org/concurrency/5
	for {
		select {
		case <-db.Shutdown: // if channel message is shutdown finish loop
			return
		case <-t.C: // on channel tick check the connection
			err := db.CheckConnection()
			if err != nil {
				log.Error(log.Database, "database connection error:", err)
			}
		}
	}
}

// Stop terminates all processes belonging to the service.
func (db *Database) Stop() error {
	if err := db.Service.Stop(); err != nil {
		return err
	}
	close(db.Shutdown)
	return nil
}

// CheckConnection checks to make sure the database is connected
func (db *Database) CheckConnection() error {

	if !db.CoreConfig.Enabled {
		return errors.ErrDatabaseSupportDisabled
	}

	if db.CoreDB == nil {
		return errors.ErrNoDatabaseProvided
	}

	if err := db.CoreDB.Ping(); err != nil {
		db.CoreDB.SetConnected(false)
		return err
	}

	if !db.CoreDB.IsConnected() {
		log.Info(log.Database, "CoreDB connection reestablished")
		db.CoreDB.SetConnected(true)
	}
	return nil
}
