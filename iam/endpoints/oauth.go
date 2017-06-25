package endpoints

import (
	"net/http"

	"log"

	"github.com/dtop/go.demo.iam/iam/models"
	"github.com/dtop/go.ginject"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/gin-server"
	"gopkg.in/oauth2.v3"
)

// Authorize is the oauth /authorize endpoint
func Authorize(c *gin.Context) {

	sess := models.NewSession(ginject.Deps(c))

	if err := c.Bind(sess); err != nil {
		c.AbortWithStatus(400)
		return
	}

	sess.Store()

	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Redirect(http.StatusTemporaryRedirect, "/gui/login?sess="+sess.GetSessionID())
}

// RealAuthorize is the actual authorization code creation
func RealAuthorize(c *gin.Context) {

	sessid := c.DefaultQuery("sess", "")
	if sessid == "" {
		c.AbortWithStatus(500)
		return
	}

	sess := models.NewSession(ginject.Deps(c))
	if err := sess.FromSessionID(sessid); err != nil {
		log.Println(err)
		c.AbortWithStatus(500)
		return
	}

	manager := models.NewManager(ginject.Deps(c))

	server.InitServer(manager)
	server.SetAllowedGrantType(oauth2.AuthorizationCode, oauth2.Refreshing)
	server.SetAllowedResponseType(oauth2.Code, oauth2.Token)

	server.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {

		return sess.GetUserID(), nil
	})

	server.HandleAuthorizeRequest(c)
}

// Token is the oauth /token endpoint
func Token(c *gin.Context) {

	clientStorage, err := models.NewClientStorage(ginject.Deps(c))
	if err != nil {
		panic(err)
		return
	}

	manager := models.NewManager(ginject.Deps(c))

	server.InitServer(manager)
	server.SetAllowedGrantType(oauth2.AuthorizationCode, oauth2.Refreshing)
	server.SetAllowedResponseType(oauth2.Token)

	server.SetClientInfoHandler(func(r *http.Request) (clientID, clientSecret string, err error) {

		cid := r.FormValue("client_id")

		if cid == "" {

		}

		info, err := clientStorage.GetByID(cid)
		if err != nil {
			return
		}

		clientID = info.GetID()
		clientSecret = info.GetSecret()
		return
	})

	server.HandleTokenRequest(c)
}
