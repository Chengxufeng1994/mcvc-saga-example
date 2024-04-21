package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Chengxufeng1994/go-saga-example/auth-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/common/bootstrap"
	"github.com/Chengxufeng1994/go-saga-example/common/middleware"
	"github.com/gin-gonic/gin"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	gohttpmetricsmiddleware "github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewGinEngine(bootstrapCfg *bootstrap.BootstrapConfig) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(middleware.CORS())
	engine.Use(otelgin.Middleware(bootstrapCfg.Application))
	mdlw := gohttpmetricsmiddleware.New(gohttpmetricsmiddleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{
			Prefix: bootstrapCfg.Application,
		}),
	})
	engine.Use(ginmiddleware.Handler("", mdlw))

	return engine
}

type HttpServer struct {
	Application     string
	bootstrapConfig *bootstrap.BootstrapConfig
	Engine          *gin.Engine
	Router          *Router
	Srv             *http.Server
}

func New(bootstrapConfig *bootstrap.BootstrapConfig, engine *gin.Engine, router *Router) *HttpServer {
	return &HttpServer{
		Application:     bootstrapConfig.Application,
		bootstrapConfig: bootstrapConfig,
		Engine:          engine,
		Router:          router,
	}
}

func (s *HttpServer) Run() error {
	s.Router.RegisterRoutes()

	addr := fmt.Sprintf(":%d", s.bootstrapConfig.HTTP.Port)
	readTimeout, _ := time.ParseDuration(s.bootstrapConfig.HTTP.ReadTimeout)
	writeTimeout, _ := time.ParseDuration(s.bootstrapConfig.HTTP.WriteTimeout)
	idleTimeout, _ := time.ParseDuration(s.bootstrapConfig.HTTP.IdleTimeout)

	s.Srv = &http.Server{
		Addr:         addr,
		Handler:      s.Engine,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	config.ContextLogger.Infoln("http.Run listening on", s.bootstrapConfig.HTTP.Port)
	if err := s.Srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *HttpServer) GracefulShutdown(ctx context.Context) {
	config.ContextLogger.Infoln("http.GracefulShutdown")
	_ = s.Srv.Shutdown(ctx)
}
