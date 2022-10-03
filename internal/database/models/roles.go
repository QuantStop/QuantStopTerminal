package models

import (
	"database/sql"
	"github.com/quantstop/quantstopterminal/internal/database/drivers"
	"github.com/quantstop/quantstopterminal/internal/log"
)

type Role struct {
	ID   int
	Name string
}

var defaultRoles = []Role{
	{Name: "user"},
	{Name: "moderator"},
	{Name: "admin"},
}

func CreateRolesTable(db *sql.DB, driver string) error {

	log.Debugln(log.Database, "Checking for roles table ...")
	var row *sql.Row
	var table interface{}

	switch driver {
	case drivers.DBPostgreSQL:
		// todo: change to postgre
		row = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='roles' LIMIT 1")
	case drivers.DBSQLite, drivers.DBSQLite3:
		row = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='roles' LIMIT 1")
	case drivers.DBMySQL:
		// todo: change to mysql
		row = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='roles' LIMIT 1")
	default:
		row = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='roles' LIMIT 1")
	}

	// returns err if no table is round
	if err := row.Scan(&table); err != nil {
		log.Debugln(log.Database, "Checking for roles table ... Not found.")
		log.Debugln(log.Database, "Creating roles table ... ")
		usersTable := `
			create table if not exists roles
			(
				id integer primary key autoincrement,
				name varchar(255) not null,
				constraint name
					unique (name)
			);
		`
		_, err := db.Exec(usersTable)
		if err != nil {
			log.Errorf(log.Database, "Creating roles table ... Failed. Error: %v", err)
			return err // todo: custom error?
		}
		log.Debugln(log.Database, "Creating roles table ... Success!")

		if err = CreateDefaultRoles(db); err != nil {
			log.Errorf(log.Database, "Error creating default roles: %v", err)
			return err
		}
	}

	log.Debugln(log.Database, "Checking for roles table ... Found!")
	return nil
}

func CreateDefaultRoles(db *sql.DB) error {

	for _, role := range defaultRoles {
		if err := role.CreateRole(db); err != nil {
			return err
		}
	}

	return nil
}

func (r *Role) CreateRole(db *sql.DB) error {

	log.Debugln(log.Database, "Creating role ...")

	result, err := db.Exec("INSERT INTO roles (name) VALUES ($1)", r.Name)
	if err != nil {
		log.Errorf(log.Database, "could not insert row: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf(log.Database, "could not get affected rows: %v", err)
		return err
	}

	log.Debugln(log.Database, "Role created. Inserted", rowsAffected, "rows")

	return nil
}
