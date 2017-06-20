package models

import (
	"encoding/base64"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dtop/go.ginject"
	"github.com/satori/go.uuid"
	"gopkg.in/oauth2.v3"
)

type TokenGenerator struct {
	keyProvider KeyProvider
}

func NewTokenGenerator(dep ginject.Injector) *TokenGenerator {

	kp, err := NewKeyProvider(dep)
	if err != nil {
		panic(err)
	}

	return &TokenGenerator{keyProvider: kp}
}

// Token generates the actual token
func (tg *TokenGenerator) Token(data *oauth2.GenerateBasic, isGenRefresh bool) (access, refresh string, err error) {

	ti := data.TokenInfo

	claims := make(jwt.MapClaims)
	// default fields
	claims["iss"] = "DemoIAM"
	claims["aud"] = "DemoServices"
	claims["exp"] = time.Now().Add(ti.GetAccessExpiresIn()).Unix()
	claims["iat"] = time.Now().Unix()
	claims["nbf"] = ""
	claims["sub"] = ""
	claims["jti"] = data.UserID
	// own fields
	claims["sco"] = ti.GetScope()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	access, err = token.SignedString(tg.keyProvider.GetPrivateKey())
	if err != nil {
		return
	}

	if isGenRefresh {

		refresh = base64.URLEncoding.EncodeToString(uuid.NewV5(uuid.NewV4(), access).Bytes())
	}
	return
}
