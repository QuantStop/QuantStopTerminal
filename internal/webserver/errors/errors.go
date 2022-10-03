package errors

import (
	"errors"
	"net/http"
)

var (
	NilDBError = errors.New("nil database")

	BadRequestMethod = errors.New(http.StatusText(http.StatusMethodNotAllowed))
	InternalError    = errors.New(http.StatusText(http.StatusInternalServerError))

	NoJSONBody   = errors.New("unable to decode JSON")
	InvalidEmail = errors.New("invalid email")
	InvalidInput = errors.New("invalid input")

	FailedLogin          = errors.New("invalid username or password")
	AlreadyRegistered    = errors.New("an account already exists for this email")
	VerificationNotFound = errors.New("invalid verification code")
	VerificationExpired  = errors.New("verification code was already used")

	UserNotFound  = errors.New("user does not exist")
	PostNotFound  = errors.New("post does not exist")
	ResetNotFound = errors.New("invalid password reset code")

	BadCSRF           = errors.New("missing CSRF header")
	BadOrigin         = errors.New("invalid origin header")
	RouteUnauthorized = errors.New("you don't have permission to view this resource")
	RouteNotFound     = errors.New("route not found")
	ExpiredToken      = errors.New("your access token expired")
	InvalidToken      = errors.New("your access token is invalid")
)

// codeMap returns a map of errors to http status codes
func codeMap() map[error]int {
	return map[error]int{
		BadRequestMethod: http.StatusMethodNotAllowed,
		InternalError:    http.StatusInternalServerError,

		NoJSONBody:        http.StatusBadRequest,
		InvalidEmail:      http.StatusBadRequest,
		AlreadyRegistered: http.StatusBadRequest,

		FailedLogin:          http.StatusUnauthorized,
		VerificationNotFound: http.StatusNotFound,
		VerificationExpired:  http.StatusGone,
		UserNotFound:         http.StatusNotFound,
		PostNotFound:         http.StatusNotFound,
		ResetNotFound:        http.StatusNotFound,

		BadCSRF:           http.StatusUnauthorized,
		BadOrigin:         http.StatusUnauthorized,
		RouteUnauthorized: http.StatusUnauthorized,
		RouteNotFound:     http.StatusNotFound,
		ExpiredToken:      http.StatusUnauthorized,
		InvalidToken:      http.StatusUnauthorized,
	}
}

// GetCode is a helper to get the relevant code for an error, or just return 500
func GetCode(e error) (bool, int) {
	if code, ok := codeMap()[e]; ok {
		return true, code
	}
	return false, http.StatusInternalServerError
}
