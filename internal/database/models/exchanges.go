package models

import (
	"database/sql"
	"errors"
	"github.com/quantstop/quantstopexchange/qsx"
	"github.com/quantstop/quantstopterminal/internal/log"
)

type Exchange struct {
	ID             int
	Name           string
	AuthKey        string
	AuthPassphrase string
	AuthSecret     string
	Currency       string
}

func CreateExchangesTable(db *sql.DB, driver string) error {

	log.Debugln(log.DatabaseLogger, "Checking for exchanges table ...")
	row := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='exchanges' LIMIT 1")
	var table interface{}

	// returns err if no table is round
	if err := row.Scan(&table); err != nil {
		log.Debugln(log.DatabaseLogger, "Checking for exchanges table ... Not found.")
		log.Debugln(log.DatabaseLogger, "Creating exchanges table ... ")
		exchangeTable := `
			create table if not exists exchanges
			(
				id integer primary key autoincrement,
				name varchar(255) not null,
				authKey varchar(255),
				authPassphrase varchar(255),
				authSecret varchar(255),
				currency varchar(255),
				constraint name
					unique (name)
			);
		`
		_, err := db.Exec(exchangeTable)
		if err != nil {
			log.Errorf(log.DatabaseLogger, "Creating exchanges table ... Failed. Error: %v", err)
			return err // todo: custom error?
		}
		log.Debugln(log.DatabaseLogger, "Creating exchanges table ... Success!")

	}

	log.Debugln(log.DatabaseLogger, "Checking for exchanges table ... Found!")
	return nil
}

func (c *Exchange) CreateExchange(db *sql.DB) error {

	log.Debugln(log.DatabaseLogger, "Creating exchange ...")

	result, err := db.Exec("INSERT INTO exchanges (name, authKey, authPassphrase, authSecret, currency) VALUES ($1, $2, $3, $4, $5)", c.Name, c.AuthKey, c.AuthPassphrase, c.AuthSecret, c.Currency)
	if err != nil {
		log.Errorf(log.DatabaseLogger, "could not insert row: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf(log.DatabaseLogger, "could not get affected rows: %v", err)
		return err
	}

	log.Debugln(log.DatabaseLogger, "exchange created. Inserted", rowsAffected, "rows")

	return nil
}

func (c *Exchange) GetExchangeByName(db *sql.DB, name qsx.Name) error {

	if name == "" {
		log.Errorf(log.DatabaseLogger, "username is nil")
		return errors.New("exchange model, cannot GetExchangeByName, username is nil")
	}

	if db == nil {
		log.Errorf(log.DatabaseLogger, "db is nil")
		return errors.New("exchange model, cannot GetExchangeByName, db is nil")
	}

	query := `
		SELECT u.id, u.name, u.authKey, u.authPassphrase, u.authSecret, u.currency
		FROM exchanges AS u
		WHERE u.name = ?
	`
	rows, err := db.Query(query, name)
	if err != nil {
		log.Errorf(log.DatabaseLogger, "error getting exchange: %v", err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(
			&c.ID,
			&c.Name,
			&c.AuthKey,
			&c.AuthPassphrase,
			&c.AuthSecret,
			&c.Currency,
		)
		if err != nil {
			log.Errorf(log.DatabaseLogger, "error scanning rows: %v", err)
			return err
		}
	}

	return nil
}
