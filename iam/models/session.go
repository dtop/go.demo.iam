package models

import (
	"encoding/json"
	"fmt"

	"math/rand"
	"time"

	"github.com/dtop/go.ginject"
	"github.com/google/go-querystring/query"
	"gopkg.in/redis.v4"
)

type (

	// Session interface
	Session interface {
		New() Session
		FromSessionID(sessID string) error
		Store(sessID ...string) error
		GetSessionID() string
		GetUserID() string
		AssignUserID(userID string)
		Assemble() (string, error)
	}

	// sess is the Session implementation
	sess struct {
		SessID       string `json:"-" url:"-"`
		UserID       string `json:"user_id" url:"-"`
		ClientID     string `json:"client_id" url:"client_id" form:"client_id" binding:"required"`
		RedirectURI  string `json:"redirect_uri" url:"redirect_uri" form:"redirect_uri" binding:"required"`
		Scope        string `json:"scope" url:"scope" form:"scope"`
		State        string `json:"state" url:"state" form:"state"`
		ResponseType string `json:"response_type" url:"response_type" form:"response_type" binding:"required"`

		deps  ginject.Injector `json:"-" url:"-"`
		Redis *redis.Client    `inject:"redis" json:"-" url:"-"`
	}
)

// NewSession creates a new session item
func NewSession(deps ginject.Injector) Session {

	sess := &sess{deps: deps}
	if err := deps.Apply(sess); err != nil {
		panic(err)
	}

	return sess
}

// ################## Session

// New creates a brand new session item
func (s *sess) New() Session {
	return NewSession(s.deps)
}

// FromSessionID loads a stored session and applies it on the object
func (s *sess) FromSessionID(sessID string) error {

	raw, err := s.Redis.Get(fmt.Sprintf("sess_%v", sessID)).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(raw), s); err != nil {
		return err
	}

	s.SessID = sessID
	return nil
}

// Store stores the session
func (s *sess) Store(sessID ...string) error {

	var _sessid string

	if len(sessID) > 0 {
		_sessid = sessID[0]
	}

	if _sessid == "" && s.SessID != "" {
		_sessid = s.SessID
	}

	if _sessid == "" {
		_sessid = randStr(32)
	}

	raw, err := json.Marshal(s)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("sess_%v", _sessid)
	val := string(raw)
	dur := 7200 * time.Second

	if _, err := s.Redis.Set(key, val, dur).Result(); err != nil {
		return err
	}

	s.SessID = _sessid
	return nil
}

// GetSessionID returns the session id
func (s *sess) GetSessionID() string {
	return s.SessID
}

// GetUserID returns the user ID
func (s *sess) GetUserID() string {
	return s.UserID
}

func (s *sess) AssignUserID(userID string) {
	s.UserID = userID
}

func (s *sess) Assemble() (string, error) {

	v, err := query.Values(s)
	if err != nil {
		return "", err
	}

	return v.Encode(), nil
}

// #################### Helpers

func randStr(length int) string {

	rand.Seed(time.Now().UnixNano())
	abc := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)
	for i := range b {
		b[i] = abc[rand.Intn(len(abc))]
	}

	return string(b)
}
