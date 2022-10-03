package postgres

import (
	"database/sql"
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/database/drivers"
	"github.com/quantstop/quantstopterminal/internal/database/errors"
	"github.com/quantstop/quantstopterminal/internal/database/repository"

	// import go libpq driver package
	_ "github.com/lib/pq"
)

// Connect opens a connection to Postgres database and returns a pointer to database.CoreDB
func Connect(name string, cfg *drivers.ConnectionDetails) (*repository.Instance, error) {
	if cfg == nil {
		return nil, errors.ErrNilConfig
	}

	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}

	configDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode)

	db, err := sql.Open("postgres", configDSN)
	if err != nil {
		return nil, err
	}
	switch name {
	case "core":
		err = repository.CoreDB.SetPostgresConnection(db)
		if err != nil {
			return nil, err
		}
		return repository.CoreDB, nil
	case "coinbase":
		err = repository.CoinbaseDB.SetPostgresConnection(db)
		if err != nil {
			return nil, err
		}
		return repository.CoinbaseDB, nil
	case "tdameritrade":
		err = repository.TDAmeritradeDB.SetPostgresConnection(db)
		if err != nil {
			return nil, err
		}
		return repository.TDAmeritradeDB, nil
	default:
		err = repository.CoreDB.SetPostgresConnection(db)
		if err != nil {
			return nil, err
		}
		return repository.CoreDB, nil
	}
}
