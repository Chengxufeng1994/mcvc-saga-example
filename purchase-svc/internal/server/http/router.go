package http

import (
	"net/http"

	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/application"
	v1 "github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/server/http/controller/v1"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/server/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Router struct {
	app    *application.Application
	engine *gin.Engine
	logger *config.Logger
}

func NewRouter(logger *config.Logger, engine *gin.Engine, app *application.Application) *Router {
	return &Router{
		app:    app,
		engine: engine,
		logger: logger,
	}
}

func (r *Router) RegisterRoutes() {
	// K8s probe for kubernetes health checks.
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "The server is up and running.")
	})

	// prometheus probe for prometheus pull;
	r.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Handling a page not found endpoint -.
	r.engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "The requested page is not found. Please try later!"})
	})

	jwtAuthenticator := middleware.NewJwtAuthenticator(r.logger, r.app.AuthService)
	purchaseController := v1.NewPurchaseController(r.app.PurchaseService)
	v1Group := r.engine.Group("/api/v1")
	purchaseGroup := v1Group.Group("/purchase")
	purchaseGroup.Use(jwtAuthenticator.Auth())
	{
		purchaseGroup.POST("", purchaseController.CreatePurchase)
		purchaseGroup.GET("/result", purchaseController.GetResult)
	}
}
