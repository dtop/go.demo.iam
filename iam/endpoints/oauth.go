package endpoints

import (
	"net/http"

	"github.com/dtop/go.demo.iam/iam/models"
	"github.com/dtop/go.ginject"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/gin-server"
	"gopkg.in/oauth2.v3"
)

func Authorize(c *gin.Context) {

	manager := models.NewManager(ginject.Deps(c))

	server.InitServer(manager)
	server.SetAllowedGrantType(oauth2.AuthorizationCode, oauth2.Refreshing)
	server.SetAllowedResponseType(oauth2.Code, oauth2.Token)

	server.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {

		// auth users here
		return "", nil
	})

	server.HandleAuthorizeRequest(c)
}

func Token(c *gin.Context) {

	manager := models.NewManager(ginject.Deps(c))

	server.InitServer(manager)
	server.SetAllowedGrantType(oauth2.AuthorizationCode, oauth2.Refreshing)
	server.SetAllowedResponseType(oauth2.Token)

	server.SetClientInfoHandler(func(r *http.Request) (clientID, clientSecret string, err error) {

		// return client infos here
		return "", "", nil
	})

	server.HandleTokenRequest(c)
}
