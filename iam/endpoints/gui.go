package endpoints

import (
	"log"

	"fmt"

	"github.com/dtop/go.demo.iam/iam/models"
	"github.com/dtop/go.ginject"
	"github.com/gin-gonic/gin"
)

// Login displays the login form
func Login(c *gin.Context) {

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

	c.HTML(200, "login.tmpl", gin.H{
		"sessid": sess.GetSessionID(),
	})
}

// ProcessLogin logs the user in (json call)
func ProcessLogin(c *gin.Context) {

	userLogin := models.NewUserLogin(ginject.Deps(c))

	if err := c.Bind(userLogin); err != nil {
		c.JSON(400, struct {
			Error        bool   `json:"error"`
			ErrorMessage string `json:"error_message"`
		}{
			Error:        true,
			ErrorMessage: err.Error(),
		})
		return
	}

	errMap := userLogin.CheckLogin()
	if len(errMap) > 0 {
		c.JSON(400, struct {
			Error  bool         `json:"error"`
			Errors []models.Err `json:"error_messages"`
		}{
			Error:  true,
			Errors: errMap,
		})
		return
	}

	sess := models.NewSession(ginject.Deps(c))
	if err := sess.FromSessionID(userLogin.SessionID); err != nil {
		log.Println(err)
		c.AbortWithStatus(500)
		return
	}

	sess.AssignUserID(userLogin.UserID)
	sess.Store()

	qry, err := sess.Assemble()
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(500)
		return
	}

	c.JSON(200, struct {
		Error    bool   `json:"error"`
		Redirect string `json:"redirect"`
	}{
		Error:    false,
		Redirect: fmt.Sprintf("/oauth/forward?sess=%v&%v", sess.GetSessionID(), qry),
	})
}
