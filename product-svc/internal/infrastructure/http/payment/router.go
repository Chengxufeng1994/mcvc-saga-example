package product

import (
	"net/http"

	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/application"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/http/middleware"
	v1 "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/http/payment/controller/v1"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Router struct {
	app              *application.PaymentApplication
	jwtAuthenticator *middleware.JwtAuthenticator
	engine           *gin.Engine
}

func NewRouter(engine *gin.Engine, app *application.PaymentApplication, jwtAuthenticator *middleware.JwtAuthenticator) *Router {
	return &Router{
		app:              app,
		jwtAuthenticator: jwtAuthenticator,
		engine:           engine,
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

	paymentController := v1.NewPaymentController(r.app.PaymentService)
	v1 := r.engine.Group("/api/v1")
	paymentGroup := v1.Group("/payment")
	paymentGroup.Use(r.jwtAuthenticator.Auth())
	{
		paymentGroup.GET("/:id", paymentController.GetPayment)
	}
}
