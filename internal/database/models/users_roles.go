package models

import (
	"database/sql"
	"github.com/quantstop/quantstopterminal/internal/log"
)

type UserRole struct {
	UserID int
	RoleID int
}

func CreateUsersRolesTable(db *sql.DB, driver string) error {

	// todo: still only sqlite, dont like this too much as it is. could do a switch/case here with driver string parm ...

	log.Debugln(log.Database, "Checking for users_roles table ...")
	row := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users_roles' LIMIT 1")
	var table interface{}

	// returns err if no table is round
	if err := row.Scan(&table); err != nil {
		log.Debugln(log.Database, "Checking for users_roles table ... Not found.")
		log.Debugln(log.Database, "Creating users_roles table ... ")
		usersTable := `
			create table if not exists users_roles
			(
				user_id integer,
				role_id integer,
				foreign key (user_id) references users(id),
				foreign key (role_id) references roles(id)
			);
		`
		_, err := db.Exec(usersTable)
		if err != nil {
			log.Errorf(log.Database, "Creating users_roles table ... Failed. Error: %v", err)
			return err // todo: custom error?
		}
		log.Debugln(log.Database, "Creating users_roles table ... Success!")

	}

	log.Debugln(log.Database, "Checking for users_roles table ... Found!")
	return nil
}

func CreateDefaultAdminRoles(db *sql.DB) error {

	for i := range defaultRoles {

		if i != 0 {
			ur := UserRole{RoleID: i + 1, UserID: 1}
			if err := ur.CreateUserRole(db); err != nil {
				return err
			}
		}

	}

	return nil
}

func (ur *UserRole) CreateUserRole(db *sql.DB) error {
	log.Debugln(log.Database, "Creating user role association ...")

	result, err := db.Exec("INSERT INTO users_roles (user_id, role_id) VALUES ($1, $2)", ur.UserID, ur.RoleID)
	if err != nil {
		log.Errorf(log.Database, "could not insert row: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf(log.Database, "could not get affected rows: %v", err)
		return err
	}

	log.Debugln(log.Database, "User role association created. Inserted", rowsAffected, "rows")

	return nil
}

/*func (ur *UserRole) GetUsersRoles(db *sql.DB) ([]Role, error) {

}*/
