package models

import (
	"encoding/base64"
	"time"

	"fmt"

	"database/sql"

	"encoding/json"

	"log"

	"github.com/dtop/go.demo.iam/iam/wrappers"
	"github.com/dtop/go.ginject"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/redis.v4"
)

const (

	// PrefixToken is the prefix for the token in redis
	PrefixToken = "token"
	// PrefixCode is the prefix for the code in redis
	PrefixCode = "code"
	// PrefixAccessToken is the prefix for the access token in redis
	PrefixAccessToken = "access_token"
)

// TokenStorage implements oauth2.TokenStore
type TokenStorage struct {
	dep   ginject.Injector
	Db    *wrappers.MySQL `inject:"db"`
	Redis *redis.Client   `inject:"redis"`
}

// NewTokenStorage creates and returns a new token storage including all dependencies set
func NewTokenStorage(dep ginject.Injector) (oauth2.TokenStore, error) {

	ts := &TokenStorage{dep: dep}
	if err := dep.Apply(ts); err != nil {
		return nil, err
	}

	return ts, nil
}

// Create creates all records of tokens after issuing
func (m TokenStorage) Create(info oauth2.TokenInfo) error {

	bin, err := json.Marshal(info)
	if err != nil {
		return err
	}

	now := time.Now()
	enc := base64.StdEncoding.EncodeToString(bin)

	aexp := info.GetAccessExpiresIn()
	rexp := aexp

	// if code is set, write only code to store
	if code := info.GetCode(); code != "" {

		if _, err := m.Redis.Set(prefixedKey(PrefixCode, code), enc, info.GetCodeExpiresIn()).Result(); err != nil {
			return err
		}

		return nil
	}

	// otherwise write rest
	if refresh := info.GetRefresh(); refresh != "" {
		rexp = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).Sub(now)
		if aexp.Seconds() > rexp.Seconds() {
			aexp = rexp
		}

		_, err := m.Db.Exec(func(db *sql.DB) (sql.Result, error) {

			qry := "REPLACE INTO refresh_tokens (code, token, expiry) VALUES (?,?,?)"
			return db.Exec(qry, refresh, enc, time.Now().Add(rexp).Format("2006-01-02 15:04:05"))
		})

		if err != nil {
			return err
		}
	}

	if _, err := m.Redis.Set(prefixedKey(PrefixAccessToken, info.GetAccess()), enc, aexp).Result(); err != nil {
		return err
	}

	return nil
}

func (m TokenStorage) getData(key string, isCode bool) (oauth2.TokenInfo, error) {

	prefix := PrefixToken
	if isCode {
		prefix = PrefixCode
	}

	enc, err := m.Redis.Get(prefixedKey(prefix, key)).Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	token, err := m.buildToken(enc, isCode)

	return token, nil
}

func (m TokenStorage) buildToken(enc string, isCode bool) (oauth2.TokenInfo, error) {

	raw, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	token := &models.Token{}
	if err := json.Unmarshal(raw, token); err != nil {
		log.Println(err)
		return nil, err
	}

	return token, nil
}

func (m TokenStorage) remove(key, prefix string) (err error) {

	_, err = m.Redis.Del(prefixedKey(prefix, key)).Result()
	return
}

// GetByCode uses the authorization code for token information data
func (m TokenStorage) GetByCode(code string) (ti oauth2.TokenInfo, err error) {

	ti, err = m.getData(code, true)
	return
}

// GetByAccess uses the access token for token information data
func (m TokenStorage) GetByAccess(access string) (ti oauth2.TokenInfo, err error) {

	ti, err = m.getData(access, false)
	return
}

// GetByRefresh uses the refresh token for token information data
func (m TokenStorage) GetByRefresh(refresh string) (ti oauth2.TokenInfo, err error) {

	rows, err := m.Db.Query(func(db *sql.DB) (*sql.Rows, error) {

		qry := "SELECT token FROM refresh_tokens WHERE code = ?"
		return db.Query(qry, refresh)
	})

	if err != nil {
		return nil, err
	}

	rows.Next()
	var token string
	rows.Scan(&token)

	return m.buildToken(token, false)
}

// RemoveByCode deletes the authorization code
func (m TokenStorage) RemoveByCode(code string) error {

	return m.remove(code, PrefixCode)
}

// RemoveByAccess uses the access token to delete the token information
func (m TokenStorage) RemoveByAccess(access string) error {

	return m.remove(access, PrefixAccessToken)
}

// RemoveByRefresh uses the refresh token to delete the token information
func (m TokenStorage) RemoveByRefresh(refresh string) error {

	_, err := m.Db.Exec(func(db *sql.DB) (sql.Result, error) {

		qry := "DELETE FROM refresh_tokens WHERE code = ?"
		return db.Exec(qry, refresh)
	})

	return err
}

func prefixedKey(prefix, key string) string {
	return fmt.Sprintf("%v_%v", prefix, key)
}
