package iam

import (
	"fmt"
	"net/http"

	"os"
	"path/filepath"

	"github.com/dtop/go.demo.iam/iam/endpoints"
	"github.com/dtop/go.demo.iam/iam/wrappers"
	"github.com/dtop/go.ginject"
	"github.com/gin-gonic/gin"
)

type (
	// Iam is the actual server object
	Iam struct {
		gin *gin.Engine
		dep ginject.Injector
	}
)

// ############### Iam

// New creates a new Iam object
func New() *Iam {

	_iam := &Iam{}
	_iam.Bootstrap()

	return _iam
}

// Bootstrap setups the service
func (iam *Iam) Bootstrap() {

	// create gin and dependency manager
	iam.createGin()

	// setup available routes
	iam.setupRoutes()

	// setup dependencies
	iam.setupDeps()

	// setup templates
	iam.setupTemplates()
}

// Run runs the service
func (iam Iam) Run() {

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", 9020),
		Handler: iam.gin,
	}

	server.ListenAndServe()
}

// ################## helpers

func (iam *Iam) createGin() {

	iam.gin = gin.New()
	iam.dep = ginject.New()

	iam.gin.Use(gin.ErrorLogger())
	iam.gin.Use(gin.Recovery())
	iam.gin.Use(ginject.DependencyInjector(iam.dep))
}

func (iam *Iam) setupRoutes() {

	oauthGroup := iam.gin.Group("/oauth")
	{
		oauthGroup.GET("/authorize", endpoints.Authorize)
		oauthGroup.POST("/token", endpoints.Token)
		oauthGroup.GET("/forward", endpoints.RealAuthorize)
	}

	guigroup := iam.gin.Group("/gui")
	{
		guigroup.GET("/login", endpoints.Login)
		guigroup.POST("/check", endpoints.ProcessLogin)
	}

	iam.gin.GET("/iam/.well-known/key", endpoints.Key)
}

func (iam *Iam) setupDeps() {

	deps := iam.dep

	deps.Register(ginject.IService(
		"db",
		wrappers.NewMySQL("", "localhost", -1, "demouser", "demopass", "demoiam"),
	))

	deps.Register(ginject.IService(
		"redis",
		wrappers.NewRedis("localhost", -1),
	))
}

func (iam *Iam) setupTemplates() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	iam.gin.LoadHTMLGlob(fmt.Sprintf("%v/templates/*", dir))
}
