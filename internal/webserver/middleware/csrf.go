package middleware

import (
	"github.com/quantstop/quantstopterminal/internal/webserver/errors"
	"github.com/quantstop/quantstopterminal/internal/webserver/write"
	"net/http"
)

// Csrf checks for the CSRF prevention header and compares the origin header
func Csrf(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//log.Println("Middleware chain | 5 | Csrf")
		if skipCorsAndCSRF(r.URL.Path) {
			fn(w, r)
			return
		}

		if r.Method != http.MethodOptions {
			if r.Header.Get("Origin") != "" && validateOrigin(r) == "" {
				// if an origin is provided, but didn't match our list
				fn = write.Error(errors.BadOrigin)
			} else if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
				fn = write.Error(errors.BadCSRF)
			}
		}
		fn(w, r)
	}
}
