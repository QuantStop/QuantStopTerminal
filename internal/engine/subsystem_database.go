package engine

import (
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/database"
	"github.com/quantstop/quantstopterminal/internal/database/drivers/mysql"
	pgsql "github.com/quantstop/quantstopterminal/internal/database/drivers/postgres"
	sqlite "github.com/quantstop/quantstopterminal/internal/database/drivers/sqlite3"
	"github.com/quantstop/quantstopterminal/internal/database/seed"
	"github.com/quantstop/quantstopterminal/internal/log"
	"sync"
	"time"
)

type DatabaseSubsystem struct {
	Subsystem
	wg                   sync.WaitGroup
	coreDatabase         *database.Instance
	coinbaseDatabase     *database.Instance
	tdameritradeDatabase *database.Instance
}

// init sets config and params
func (s *DatabaseSubsystem) init(bot *Engine, name string) error {
	if err := s.Subsystem.init(bot, name, true); err != nil {
		return err
	}
	s.coreDatabase = database.CoreDB
	if err := s.coreDatabase.SetConfig(&s.bot.Config.CoreDB); err != nil {
		return err
	}
	s.coinbaseDatabase = database.CoinbaseDB
	if err := s.coinbaseDatabase.SetConfig(&s.bot.Config.CoinbaseDB); err != nil {
		return err
	}
	s.tdameritradeDatabase = database.TDAmeritradeDB
	if err := s.tdameritradeDatabase.SetConfig(&s.bot.Config.TDAmeritradeDB); err != nil {
		return err
	}
	s.initialized = true
	log.Debugln(log.DatabaseLogger, s.name+MsgSubsystemInitialized)
	return nil
}

// start sets up the database subsystem to maintain an SQL connection
func (s *DatabaseSubsystem) start(wg *sync.WaitGroup) (err error) {
	if err = s.Subsystem.start(wg); err != nil {
		return err
	}

	if s.bot.Config.CoreDB.Enabled {
		switch s.bot.Config.CoreDB.Driver {
		case database.DBPostgreSQL:
			log.Debugf(log.DatabaseLogger,
				"database subsystem attempting to establish database connection to host %s/%s utilising %s driver\n",
				s.bot.Config.CoreDB.DSN.Host,
				s.bot.Config.CoreDB.DSN.Database,
				s.bot.Config.CoreDB.Driver)
			s.coreDatabase, err = pgsql.Connect("core", &s.bot.Config.CoreDB.DSN)
		case database.DBSQLite, database.DBSQLite3:
			log.Debugf(log.DatabaseLogger,
				"database subsystem attempting to establish database connection to %s utilising %s driver\n",
				s.bot.Config.CoreDB.DSN.Database,
				s.bot.Config.CoreDB.Driver)
			s.coreDatabase, err = sqlite.Connect("core", s.bot.Config.CoreDB.DSN.Database)
		case database.DBMySQL:
			log.Debugf(log.DatabaseLogger,
				"database subsystem attempting to establish database connection to host %s/%s utilising %s driver\n",
				s.bot.Config.CoreDB.DSN.Host,
				s.bot.Config.CoreDB.DSN.Database,
				s.bot.Config.CoreDB.Driver)
			s.coreDatabase, err = mysql.Connect("core", &s.bot.Config.CoreDB.DSN)
		default:
			return database.ErrNoDatabaseProvided
		}
		if err != nil {
			return fmt.Errorf("%w: %v Core database unavailable", database.ErrFailedToConnect, err)
		}
		s.coreDatabase.SetConnected(true)
	}

	if s.bot.Config.CoinbaseDB.Enabled {
		switch s.bot.Config.CoinbaseDB.Driver {
		case database.DBPostgreSQL:
			log.Debugf(log.DatabaseLogger,
				"database subsystem attempting to establish database connection to host %s/%s utilising %s driver\n",
				s.bot.Config.CoinbaseDB.DSN.Host,
				s.bot.Config.CoinbaseDB.DSN.Database,
				s.bot.Config.CoinbaseDB.Driver)
			s.coinbaseDatabase, err = pgsql.Connect("coinbase", &s.bot.Config.CoinbaseDB.DSN)
		case database.DBSQLite, database.DBSQLite3:
			log.Debugf(log.DatabaseLogger,
				"database subsystem attempting to establish database connection to %s utilising %s driver\n",
				s.bot.Config.CoinbaseDB.DSN.Database,
				s.bot.Config.CoinbaseDB.Driver)
			s.coinbaseDatabase, err = sqlite.Connect("coinbase", s.bot.Config.CoinbaseDB.DSN.Database)
		case database.DBMySQL:
			log.Debugf(log.DatabaseLogger,
				"database subsystem attempting to establish database connection to host %s/%s utilising %s driver\n",
				s.bot.Config.CoinbaseDB.DSN.Host,
				s.bot.Config.CoinbaseDB.DSN.Database,
				s.bot.Config.CoinbaseDB.Driver)
			s.coinbaseDatabase, err = mysql.Connect("coinbase", &s.bot.Config.CoinbaseDB.DSN)
		default:
			return database.ErrNoDatabaseProvided
		}
		if err != nil {
			return fmt.Errorf("%w: %v Coinbase database unavailable", database.ErrFailedToConnect, err)
		}
		s.coinbaseDatabase.SetConnected(true)
	}

	if s.bot.Config.TDAmeritradeDB.Enabled {
		switch s.bot.Config.TDAmeritradeDB.Driver {
		case database.DBPostgreSQL:
			log.Debugf(log.DatabaseLogger,
				"database subsystem attempting to establish database connection to host %s/%s utilising %s driver\n",
				s.bot.Config.TDAmeritradeDB.DSN.Host,
				s.bot.Config.TDAmeritradeDB.DSN.Database,
				s.bot.Config.TDAmeritradeDB.Driver)
			s.tdameritradeDatabase, err = pgsql.Connect("tdameritrade", &s.bot.Config.TDAmeritradeDB.DSN)
		case database.DBSQLite, database.DBSQLite3:
			log.Debugf(log.DatabaseLogger,
				"database subsystem attempting to establish database connection to %s utilising %s driver\n",
				s.bot.Config.TDAmeritradeDB.DSN.Database,
				s.bot.Config.TDAmeritradeDB.Driver)
			s.tdameritradeDatabase, err = sqlite.Connect("tdameritrade", s.bot.Config.TDAmeritradeDB.DSN.Database)
		case database.DBMySQL:
			log.Debugf(log.DatabaseLogger,
				"database subsystem attempting to establish database connection to host %s/%s utilising %s driver\n",
				s.bot.Config.TDAmeritradeDB.DSN.Host,
				s.bot.Config.TDAmeritradeDB.DSN.Database,
				s.bot.Config.TDAmeritradeDB.Driver)
			s.tdameritradeDatabase, err = mysql.Connect("tdameritrade", &s.bot.Config.TDAmeritradeDB.DSN)
		default:
			return database.ErrNoDatabaseProvided
		}
		if err != nil {
			return fmt.Errorf("%w: %v TD-Ameritrade database unavailable", database.ErrFailedToConnect, err)
		}
		s.tdameritradeDatabase.SetConnected(true)
	}

	// finished and connected, try seeding databases
	log.Debugln(log.DatabaseLogger, s.name+MsgSubsystemStarted)
	s.started = true

	// todo: seed all databases
	err = seed.DatabaseSeed(s.coreDatabase.SQL, s.bot.Config.CoreDB.Driver)
	if err != nil {
		return err
	}

	wg.Add(1)
	s.wg.Add(1)
	go s.run(wg)
	return nil
}

// stop attempts to shut down the subsystem
func (s *DatabaseSubsystem) stop() error {
	if err := s.Subsystem.stop(); err != nil {
		return err
	}

	s.started = false
	err := s.coreDatabase.CloseConnection()
	if err != nil {
		log.Errorf(log.DatabaseLogger, "Failed to close database: %v", err)
	}

	close(s.shutdown)
	s.wg.Wait()
	log.Debugln(log.DatabaseLogger, s.name+MsgSubsystemShutdown)
	return nil
}

// run this is the main loop for the subsystem
func (s *DatabaseSubsystem) run(wg *sync.WaitGroup) {

	t := time.NewTicker(time.Second * 5)

	defer func() {
		t.Stop()
		s.wg.Done()
		wg.Done()
		log.Debugln(log.DatabaseLogger, "CoreDB subsystem shutdown.")
	}()

	// This lets the goroutine wait on communication from the channel
	// Docs: https://tour.golang.org/concurrency/5
	for {
		select {
		case <-s.shutdown: // if channel message is shutdown finish loop
			return
		case <-t.C: // on channel tick check the connection
			err := s.CheckConnection()
			if err != nil {
				log.Error(log.DatabaseLogger, "CoreDB connection error:", err)
			}
		}
	}
}

// GetInstance returns a limited scoped database instance
func (s *DatabaseSubsystem) GetInstance() database.IDatabase {
	if s == nil || !s.started {
		return nil
	}
	return s.coreDatabase
}

// CheckConnection checks to make sure the database is connected
func (s *DatabaseSubsystem) CheckConnection() error {
	if s == nil {
		return fmt.Errorf("%s %w", "DatabaseSubsystem", ErrNilSubsystem)
	}
	if s.started == false {
		return fmt.Errorf("%s %w", "DatabaseSubsystem", ErrSubsystemNotStarted)
	}
	if !s.bot.Config.CoreDB.Enabled {
		return database.ErrDatabaseSupportDisabled
	}
	if s.coreDatabase == nil {
		return database.ErrNoDatabaseProvided
	}

	if err := s.coreDatabase.Ping(); err != nil {
		s.coreDatabase.SetConnected(false)
		return err
	}

	if !s.coreDatabase.IsConnected() {
		log.Info(log.DatabaseLogger, "CoreDB connection reestablished")
		s.coreDatabase.SetConnected(true)
	}
	return nil
}
