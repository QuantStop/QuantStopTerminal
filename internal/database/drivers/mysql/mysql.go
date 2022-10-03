package mysql

import (
	"database/sql"
	"fmt"
	"github.com/quantstop/quantstopterminal/internal/database/drivers"
	"github.com/quantstop/quantstopterminal/internal/database/errors"
	"github.com/quantstop/quantstopterminal/internal/database/repository"

	_ "github.com/go-sql-driver/mysql"
)

// Connect opens a connection to MySQL database and returns a pointer to repository.Instance
func Connect(name string, cfg *drivers.ConnectionDetails) (*repository.Instance, error) {
	if cfg == nil {
		return nil, errors.ErrNilConfig
	}

	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}

	configDSN := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		"charset=utf8mb4&parseTime=True&loc=Local")

	db, err := sql.Open("mysql", configDSN)
	if err != nil {
		return nil, err
	}

	switch name {
	case "core":
		err = repository.CoreDB.SetMySQLConnection(db)
		if err != nil {
			return nil, err
		}
		return repository.CoreDB, nil
	case "coinbase":
		err = repository.CoinbaseDB.SetMySQLConnection(db)
		if err != nil {
			return nil, err
		}
		return repository.CoinbaseDB, nil
	case "tdameritrade":
		err = repository.TDAmeritradeDB.SetMySQLConnection(db)
		if err != nil {
			return nil, err
		}
		return repository.TDAmeritradeDB, nil
	default:
		err = repository.CoreDB.SetMySQLConnection(db)
		if err != nil {
			return nil, err
		}
		return repository.CoreDB, nil
	}

}
