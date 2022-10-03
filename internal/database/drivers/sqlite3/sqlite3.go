package sqlite

import (
	"database/sql"
	"github.com/quantstop/quantstopterminal/internal/database/errors"
	"github.com/quantstop/quantstopterminal/internal/database/repository"

	//_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
	"path/filepath"
)

// Connect opens a connection to sqlite database and returns a pointer to repository.Instance
func Connect(name, db string) (*repository.Instance, error) {
	if db == "" {
		return nil, errors.ErrNoDatabaseProvided
	}

	databaseFullLocation := filepath.Join(repository.CoreDB.GetConfig().ConfigDir, db)
	dbConn, err := sql.Open("sqlite", databaseFullLocation)
	if err != nil {
		return nil, err
	}

	switch name {
	case "core":
		err = repository.CoreDB.SetSQLiteConnection(dbConn)
		if err != nil {
			return nil, err
		}
		return repository.CoreDB, nil
	case "coinbase":
		err = repository.CoinbaseDB.SetSQLiteConnection(dbConn)
		if err != nil {
			return nil, err
		}
		return repository.CoinbaseDB, nil
	case "tdameritrade":
		err = repository.TDAmeritradeDB.SetSQLiteConnection(dbConn)
		if err != nil {
			return nil, err
		}
		return repository.TDAmeritradeDB, nil
	default:
		err = repository.CoreDB.SetSQLiteConnection(dbConn)
		if err != nil {
			return nil, err
		}
		return repository.CoreDB, nil
	}
}
