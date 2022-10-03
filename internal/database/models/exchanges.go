package models

import (
	"database/sql"
	"errors"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
)

type Exchange struct {
	ID             int
	UserDefined    int
	Name           string
	AuthKey        string
	AuthPassphrase string
	AuthSecret     string
	Token          string
	Currency       string
}

func CreateExchangesTable(db *sql.DB, driver string) error {

	log.Debugln(log.Database, "Checking for exchanges table ...")
	row := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='exchanges' LIMIT 1")
	var table interface{}

	// returns err if no table is round
	if err := row.Scan(&table); err != nil {
		log.Debugln(log.Database, "Checking for exchanges table ... Not found.")
		log.Debugln(log.Database, "Creating exchanges table ... ")
		exchangeTable := `
			create table if not exists exchanges
			(
				id integer primary key autoincrement,
				userDefined integer,
				name varchar(255) not null,
				authKey varchar(255),
				authPassphrase varchar(255),
				authSecret varchar(255),
				token varchar(255),
				currency varchar(255),
				constraint name
					unique (name)
			);
		`
		_, err := db.Exec(exchangeTable)
		if err != nil {
			log.Errorf(log.Database, "Creating exchanges table ... Failed. Error: %v", err)
			return err // todo: custom error?
		}
		log.Debugln(log.Database, "Creating exchanges table ... Success!")

	}

	log.Debugln(log.Database, "Checking for exchanges table ... Found!")
	return nil
}

func CreateDefaultExchanges(db *sql.DB) error {
	log.Debugln(log.Database, "Checking for default exchanges ...")

	for _, name := range qsx.SupportedExchanges {
		if CheckDefaultExchangeExists(db, string(name)) {
			continue
		}
		exchange := Exchange{
			UserDefined:    0,
			Name:           string(name),
			AuthKey:        "",
			AuthPassphrase: "",
			AuthSecret:     "",
			Token:          "",
			Currency:       "usd", // todo: how to we handle this
		}
		if err := exchange.CreateExchange(db); err != nil {
			return err
		}
	}

	log.Debugln(log.Database, "Checking for default exchanges ... Finished!")
	return nil
}

func CheckDefaultExchangeExists(db *sql.DB, name string) bool {
	row := db.QueryRow("SELECT * FROM exchanges WHERE name=$1 LIMIT 1", name)
	e := &Exchange{}
	if err := row.Scan(&e.ID, &e.UserDefined, &e.Name, &e.AuthKey, &e.AuthPassphrase, &e.AuthSecret, &e.Token, &e.Currency); err != nil {
		return false
	}
	return true
}

func (c *Exchange) CreateExchange(db *sql.DB) error {

	log.Debugln(log.Database, "Creating exchange ...")

	query := `
		INSERT INTO exchanges (name, userDefined, authKey, authPassphrase, authSecret, token, currency) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	result, err := db.Exec(query, c.Name, c.UserDefined, c.AuthKey, c.AuthPassphrase, c.AuthSecret, c.Token, c.Currency)

	if err != nil {
		log.Errorf(log.Database, "could not insert row: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf(log.Database, "could not get affected rows: %v", err)
		return err
	}

	log.Debugln(log.Database, "exchange created. Inserted", rowsAffected, "rows")

	return nil
}

func (c *Exchange) GetExchangeByName(db *sql.DB, name qsx.Name) error {

	if name == "" {
		log.Errorf(log.Database, "name is nil")
		return errors.New("exchange model, cannot GetExchangeByName, name is nil")
	}

	if db == nil {
		log.Errorf(log.Database, "db is nil")
		return errors.New("exchange model, cannot GetExchangeByName, db is nil")
	}

	query := `
		SELECT u.id, u.userDefined, u.name, u.authKey, u.authPassphrase, u.authSecret, u.token, u.currency
		FROM exchanges AS u
		WHERE u.name = ?
	`
	rows, err := db.Query(query, name)
	if err != nil {
		log.Errorf(log.Database, "error getting exchange: %v", err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&c.ID,
			&c.UserDefined,
			&c.Name,
			&c.AuthKey,
			&c.AuthPassphrase,
			&c.AuthSecret,
			&c.Token,
			&c.Currency,
		)
		if err != nil {
			log.Errorf(log.Database, "error scanning rows: %v", err)
			return err
		}
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func (c *Exchange) GetAllExchanges(db *sql.DB) ([]Exchange, error) {

	if db == nil {
		log.Errorf(log.Database, "db is nil")
		return nil, errors.New("exchange model, cannot GetAllExchanges, db is nil")
	}

	rows, err := db.Query("SELECT * FROM exchanges")
	if err != nil {
		log.Errorf(log.Database, "error getting exchanges: %v", err)
		return nil, err
	}
	defer rows.Close()

	var exchanges []Exchange

	for rows.Next() {

		exchange := Exchange{}
		err = rows.Scan(
			&exchange.ID,
			&exchange.UserDefined,
			&exchange.Name,
			&exchange.AuthKey,
			&exchange.AuthPassphrase,
			&exchange.AuthSecret,
			&exchange.Token,
			&exchange.Currency,
		)
		exchanges = append(exchanges, exchange)
		if err != nil {
			log.Errorf(log.Database, "error scanning rows: %v", err)
			return nil, err
		}

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return exchanges, nil

}
