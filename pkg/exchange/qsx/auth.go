package qsx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"time"
)

// Auth is the type for storing all authentication information for exchanges that support it
// Please take notice that not all fields are required for every exchange
type Auth struct {
	Key        string
	Passphrase string
	Secret     string
	Token      *oauth2.Token
}

// NewAuth returns a pointer to a new Auth struct
func NewAuth(key string, passphrase string, secret string) *Auth {
	return &Auth{
		Key:        key,
		Passphrase: passphrase,
		Secret:     secret,
		Token: &oauth2.Token{
			AccessToken:  "",
			TokenType:    "",
			RefreshToken: "",
			Expiry:       time.Time{},
		},
	}
}

// SignSHA256HMAC creates a sha256 HMAC using the base64-decoded Secret key on the pre-hash string:
//
//	`timestamp + method + requestPath + body`
//
// where + represents string concatenation, and then base64 encoding the output
func SignSHA256HMAC(message string, secret string) (string, error) {

	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", errors.New(fmt.Sprintf("qsx auth error decoding secret: %v", err))
	}
	signature := hmac.New(sha256.New, key)
	_, err = signature.Write([]byte(message))
	if err != nil {
		return "", errors.New(fmt.Sprintf("qsx auth error writing signature: %v", err))
	}
	return base64.StdEncoding.EncodeToString(signature.Sum(nil)), nil
}
