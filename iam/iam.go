package iam

import (
	"fmt"
	"net/http"

	"github.com/dtop/go.demo.iam/iam/endpoints"
	"github.com/dtop/go.demo.iam/iam/wrappers"
	"github.com/dtop/go.ginject"
	"github.com/gin-gonic/gin"
)

type (
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
	}

	//v1Group := iam.gin.Group("/v1")
	//{
	//
	//	guiGroup := v1Group.Group("/gui")
	//	{
	//
	//	}
	//}
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
