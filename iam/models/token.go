package models

import (
	"crypto/sha1"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dtop/go.ginject"
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
		dep         ginject.Injector
		keyProvider KeyProvider
	}
)

// NewTokenGenerator returns a brand new token generator
func NewTokenGenerator(dep ginject.Injector) *JWTGen {

	kp, err := NewKeyProvider(dep)
	if err != nil {
		panic(err)
	}

	gen := &JWTGen{dep: dep, keyProvider: kp}
	return gen
}

// New returns a brand new token
func (ag *JWTGen) New() oauth2.TokenInfo {
	return NewTokenGenerator(ag.dep)
}

// Token generates the actual token
func (ag *JWTGen) Token(data *oauth2.GenerateBasic, isGenRefresh bool) (access, refresh string, err error) {

	claims := make(jwt.MapClaims)
	// default fields
	claims["iss"] = "DemoIAM"
	claims["aud"] = data.Client
	claims["exp"] = time.Now().Add(ag.GetAccessExpiresIn()).Unix()
	claims["iat"] = time.Now().Unix()
	claims["nbf"] = ""
	claims["sub"] = ""
	claims["jti"] = data.UserID
	// own fields
	claims["sco"] = ag.GetScope()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	access, err = token.SignedString(ag.keyProvider.GetPrivateKey())
	if err != nil {
		return
	}

	if isGenRefresh {

		h := sha1.New()
		io.WriteString(h, access)
		ref := fmt.Sprintf("% x", h.Sum(nil))
		refresh = strings.ToUpper(ref)
	}
	return
}
