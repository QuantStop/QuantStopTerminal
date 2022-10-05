package repository

import (
	"errors"
	"github.com/quantstop/quantstopterminal/internal/database/models"
	"github.com/quantstop/quantstopterminal/internal/log"
)

type CoreDatabase struct {
	*Instance
}

func (db *CoreDatabase) GetUserByUsername(username string) (*models.User, error) {
	if username == "" {
		log.Errorf(log.Database, "username is nil")
		return nil, errors.New("users model, cannot GetUserByUsername, username is nil")
	}

	if db == nil {
		log.Errorf(log.Database, "db is nil")
		return nil, errors.New("users model, cannot GetUserByUsername, db is nil")
	}

	query := `
		SELECT u.id, u.username, u.password, u.salt, r.name
		FROM users AS u
		JOIN users_roles AS ur ON u.id = ur.user_id
	    JOIN roles AS r ON ur.role_id = r.id
		WHERE u.username = ?
	`

	rows, err := db.SQL.Query(query, username)
	if err != nil {
		log.Errorf(log.Database, "error getting user: %v", err)
		return nil, err
	}

	u := &models.User{}

	for rows.Next() {
		roles := &models.Role{}
		err = rows.Scan(
			&u.ID,
			&u.Username,
			&u.Password,
			&u.Salt,
			&roles.Name,
		)
		if err != nil {
			log.Errorf(log.Database, "error scanning rows: %v", err)
			return nil, err
		}
		u.Roles = append(u.Roles, roles.Name)
	}

	return u, nil
}
