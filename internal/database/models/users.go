package models

import (
	"database/sql"
	"errors"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/internal/webserver/utils"
)

type User struct {
	ID       uint32
	Username string
	Password string
	Salt     string
	Roles    []string
}

func CreateUsersTable(db *sql.DB, driver string) error {

	// todo: still only sqlite, dont like this too much as it is. could do a switch/case here with driver string parm ...

	log.Debugln(log.Database, "Checking for users table ...")
	row := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users' LIMIT 1")
	var table interface{}

	// returns err if no table is round
	if err := row.Scan(&table); err != nil {
		log.Debugln(log.Database, "Checking for users table ... Not found.")
		log.Debugln(log.Database, "Creating users table ... ")
		usersTable := `
			create table if not exists users
			(
				id integer primary key autoincrement,
				username varchar(255) not null,
				password varchar(255) not null,
				salt varchar(255) not null,
				constraint username
					unique (username)
			);
		`
		_, err := db.Exec(usersTable)
		if err != nil {
			log.Errorf(log.Database, "Creating users table ... Failed. Error: %v", err)
			return err // todo: custom error?
		}
		log.Debugln(log.Database, "Creating users table ... Success!")

	}

	log.Debugln(log.Database, "Checking for users table ... Found!")
	return nil
}

func CreateDefaultAdmin(db *sql.DB) error {

	// Check if default admin exists
	log.Debugln(log.Database, "Checking if default admin exists ...")
	if CheckDefaultAdminExists(db) {
		return nil
	}

	// Create default admin
	log.Debugln(log.Database, "Creating default admin ... ")
	defaultUser := User{
		Username: "admin",
		Password: "admin",
	}

	var err error

	// Set salt, and hash password with salt
	defaultUser.Salt = utils.GenerateRandomString(32)
	defaultUser.Password, err = utils.HashPassword(defaultUser.Password, defaultUser.Salt)
	if err != nil {
		return err
	}

	err = defaultUser.CreateUser(db)
	if err != nil {
		return err
	}
	log.Debugln(log.Database, "Creating default admin ... Success! Finished SeedDB.")

	// check/create default admin roles
	if err = CreateDefaultAdminRoles(db); err != nil {
		return err
	}

	return nil
}

func CheckDefaultAdminExists(db *sql.DB) bool {
	row := db.QueryRow("SELECT * FROM users WHERE id=$1 LIMIT 1", "1")
	u := &User{}
	if err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Salt); err != nil {
		return false
	}
	return true
}

func (u *User) CreateUser(db *sql.DB) error {

	log.Debugln(log.Database, "Creating user ...")

	result, err := db.Exec("INSERT INTO users (username, password, salt) VALUES ($1, $2, $3)", u.Username, u.Password, u.Salt)
	if err != nil {
		log.Errorf(log.Database, "could not insert row: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf(log.Database, "could not get affected rows: %v", err)
		return err
	}

	log.Debugln(log.Database, "User created. Inserted", rowsAffected, "rows")

	// get id
	row := db.QueryRow("SELECT id FROM users WHERE username=$1", u.Username)
	if err = row.Scan(&u.ID); err != nil {
		return err
	}

	// set role
	res, err := db.Exec("INSERT INTO users_roles (user_id, role_id) VALUES ($1, $2)", u.ID, 1)
	if err != nil {
		log.Errorf(log.Database, "could not insert row: %v", err)
		return err
	}
	rowsA, err := res.RowsAffected()
	if err != nil {
		log.Errorf(log.Database, "could not get affected rows: %v", err)
		return err
	}
	log.Debugln(log.Database, "Role association created. Inserted", rowsA, "rows")

	return nil
}

func (u *User) GetUserByUsername(db *sql.DB, username string) error {

	if username == "" {
		log.Errorf(log.Database, "username is nil")
		return errors.New("users model, cannot GetUserByUsername, username is nil")
	}

	if db == nil {
		log.Errorf(log.Database, "db is nil")
		return errors.New("users model, cannot GetUserByUsername, db is nil")
	}

	query := `
		SELECT u.id, u.username, u.password, u.salt, r.name
		FROM users AS u
		JOIN users_roles AS ur ON u.id = ur.user_id
	    JOIN roles AS r ON ur.role_id = r.id
		WHERE u.username = ?
	`
	rows, err := db.Query(query, username)
	if err != nil {
		log.Errorf(log.Database, "error getting user: %v", err)
		return err
	}

	for rows.Next() {
		roles := &Role{}
		err = rows.Scan(
			&u.ID,
			&u.Username,
			&u.Password,
			&u.Salt,
			&roles.Name,
		)
		if err != nil {
			log.Errorf(log.Database, "error scanning rows: %v", err)
			return err
		}
		u.Roles = append(u.Roles, roles.Name)
	}

	return nil
}

func (u *User) GetUsers(db *sql.DB) ([]*User, error) {

	if db == nil {
		log.Errorf(log.Database, "db is nil")
		return nil, errors.New("users model, cannot GetUsers, db is nil")
	}

	query := `
		SELECT * FROM users
		JOIN users_roles ON users.id = users_roles.user_id
		JOIN roles ON users_roles.role_id = roles.id
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Errorf(log.Database, "error getting users: %v", err)
		return nil, err
	}

	var usersArray []*User

	tempID := 0
	for rows.Next() {

		user := &User{}
		role := &Role{}
		userRole := &UserRole{}
		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.Salt,
			&userRole.UserID,
			&userRole.RoleID,
			&role.ID,
			&role.Name,
		)
		if err != nil {
			log.Errorf(log.Database, "error scanning rows: %v", err)
			return nil, err
		}

		if int(user.ID) == tempID {
			for _, usr := range usersArray {
				if usr.ID == user.ID {
					usr.Roles = append(usr.Roles, role.Name) // second, third, ... roles
				}
			}
		} else {
			user.Roles = append(user.Roles, role.Name) // first role
			usersArray = append(usersArray, user)
		}

		tempID = int(user.ID)

	}

	// todo: pagination?
	// todo: is this the most efficient way?

	return usersArray, nil

}
