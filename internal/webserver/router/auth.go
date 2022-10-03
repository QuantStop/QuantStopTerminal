package router

import (
	"github.com/quantstop/quantstopterminal/internal/database/models"
	"github.com/quantstop/quantstopterminal/internal/webserver/errors"
	"github.com/quantstop/quantstopterminal/internal/webserver/jwt"
	"github.com/quantstop/quantstopterminal/internal/webserver/write"
	"net/http"
)

type AuthType int

const (
	Public AuthType = iota
	User
	Moderator
	Admin
)

func (s AuthType) String() string {
	switch s {
	case Public:
		return "public"
	case User:
		return "user"
	case Moderator:
		return "moderator"
	case Admin:
		return "admin"
	}
	return "unknown"
}

// AuthHandler is the extended handler function that our API routes use
type AuthHandler func(user *models.User, w http.ResponseWriter, r *http.Request) http.HandlerFunc

// AuthRoute populates the custom AuthHandler args for our route handlers
func AuthRoute(h AuthHandler, w http.ResponseWriter, r *http.Request, authType AuthType) http.HandlerFunc {

	//log.Println("AuthRoute: " + authType.String())

	var user *models.User

	// don't try auth on public routes!
	if authType != Public {

		// parse the user cookie
		var err error
		user, err = jwt.HandleUserCookie(w, r)
		if err != nil {
			return write.Error(err)
		}

		// If we get here with no error, it means the jwt is valid
		// continue with role based auth

		// check if there are any roles at all
		if len(user.Roles) == 0 {
			return write.Error(errors.RouteUnauthorized)
		}

		// set a flag for if we find a role match
		matchFound := false

		// loop through users roles, and set match if user role is the same as the routes auth type
		for _, role := range user.Roles {
			if role == authType.String() {
				matchFound = true
			}
		}

		// if no match found, user is not authorized
		if !matchFound {
			return write.Error(errors.RouteUnauthorized)
		}

	}

	return h(user, w, r)
}
