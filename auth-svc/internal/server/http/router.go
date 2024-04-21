package http

import (
	"net/http"

	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/application"
	v1 "github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/server/http/controller/v1"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/server/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Router struct {
	app    *application.Application
	engine *gin.Engine
}

func NewRouter(engine *gin.Engine, app *application.Application) *Router {
	return &Router{
		app:    app,
		engine: engine,
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

	jwtAuthenticator := middleware.NewJwtAuthenticator(r.app.AuthService)

	authController := v1.NewAuthController(r.app.AuthService)
	userController := v1.NewUserController(r.app.UserService)
	v1Group := r.engine.Group("/api/v1")
	authGroup := v1Group.Group("/auth")
	{
		authGroup.POST("/signup", authController.SignUp)
		authGroup.POST("/signin", authController.SignIn)
		authGroup.POST("/signout", jwtAuthenticator.Auth(), authController.SignOut)
		authGroup.POST("/refresh", authController.Refresh)
	}

	userGroup := v1Group.Group("/user")
	userGroup.Use(jwtAuthenticator.Auth())
	{
		userGroup.GET("/:id", userController.GetUserByID)
	}
}
