package models

import (
	"time"

	"github.com/dtop/go.ginject"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
)

// NewManager returns a new manager for the server
func NewManager(dep ginject.Injector) *manage.Manager {

	gen := NewTokenGenerator(dep)
	mgr := manage.NewManager()

	mgr.MapTokenModel(gen)
	mgr.MapAccessGenerate(gen)
	mgr.MapAuthorizeGenerate(generates.NewAuthorizeGenerate())

	mgr.MustTokenStorage(NewTokenStorage(dep))
	mgr.MustClientStorage(NewClientStorage(dep))
	//mgr.SetValidateURIHandler(validateURI)
	mgr.SetAuthorizeCodeTokenCfg(&manage.Config{
		AccessTokenExp:    time.Hour * 2,
		RefreshTokenExp:   time.Hour * 24 * 365, // make it one year refreshable
		IsGenerateRefresh: true,
	})

	return mgr
}
