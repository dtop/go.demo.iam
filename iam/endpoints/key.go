package endpoints

import (
	"log"

	"github.com/dtop/go.demo.iam/iam/models"
	"github.com/dtop/go.ginject"
	"github.com/gin-gonic/gin"
	"github.com/mendsley/gojwk"
)

// Key is the endpoint that exposes the public key
func Key(c *gin.Context) {

	deps := ginject.Deps(c)
	keyProvider, err := models.NewKeyProvider(deps)
	if err != nil {
		return
	}

	key := keyProvider.GetJWK()
	if key == nil {
		c.AbortWithStatus(500)
		return
	}

	json, err := gojwk.Marshal(key)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(500)
	}

	c.Data(200, "application/json", json)
}
