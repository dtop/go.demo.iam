package models

import (
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
)

type (
	helper struct {
		models.Token
	}

	// JWTGen is the dummy helper for the oauth2 Manager
	JWTGen struct {
		helper
	}
)

// NewToken returns a brand new token generator
func NewToken() *JWTGen {

	gen := &JWTGen{}
	return gen
}

// New returns a brand new token
func (ag *JWTGen) New() oauth2.TokenInfo {
	return NewToken()
}
