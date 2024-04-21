package main

import (
	"context"
	"fmt"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Chengxufeng1994/go-saga-example/auth-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/db"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/adapter/repository/postgres"
	reporedis "github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/adapter/repository/redis"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/application"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/server"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/server/grpc"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/server/http"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/server/observe"
	"github.com/Chengxufeng1994/go-saga-example/common/bootstrap"
	libconfig "github.com/Chengxufeng1994/go-saga-example/common/config"
	"github.com/Chengxufeng1994/go-saga-example/common/redis"
	"github.com/Chengxufeng1994/go-saga-example/common/token"
)

func main() {
	bootCfg := bootstrap.LoadBootstrapConfig("")
	appCfg := libconfig.LoadApplicationConfig("")
	config.InitLogger(appCfg.LogConfig.Level, bootCfg)
	_ = observe.NewTracer(bootCfg, appCfg)

	// database connection
	gormDb, err := db.NewDatabase(appCfg, config.GormLogger)
	if err != nil {
		config.ContextLogger.WithError(err).Fatal("database connection error")
	}
	config.ContextLogger.Infoln("database connect successfully")

	// migrate
	migrator := db.NewMigrator(gormDb)
	if err := migrator.Migrate(); err != nil {
		config.ContextLogger.WithError(err).Fatal("database migrate error")
	}
	config.ContextLogger.Infoln("database migrate successfully")

	rcc, err := redis.NewClusterClient(appCfg)
	if err != nil {
		config.ContextLogger.WithError(err).Fatal("redis connection error")
	}
	config.ContextLogger.Infoln("redis connect successfully")
	tokenEnhancer := token.NewJWTEnhancer([]byte(appCfg.JWTConfig.Secret))
	// initialize token repository
	tokenRepository := reporedis.NewTokenRepository(rcc)
	// initialize user repository
	userRepository := postgres.NewUserRepository(gormDb)
	// initialize auth service
	authService := application.NewAuthService(
		userRepository,
		tokenRepository,
		tokenEnhancer,
		appCfg.JWTConfig.AccessTokenExpires,
		appCfg.JWTConfig.RefreshTokenExpires)
	// initialize user service
	userService := application.NewUserService(userRepository)
	// initialize application
	app := application.New(authService, userService)
	// initialize gin engine
	engine := http.NewGinEngine(bootCfg)
	// initialize route
	router := http.NewRouter(engine, app)
	// initialize http server
	httpSrv := http.New(bootCfg, engine, router)
	// initialize grpc server
	grpcSrv := grpc.New(bootCfg, authService)
	// initialize server
	srv := server.New(httpSrv, grpcSrv)

	go func() {
		if err := srv.Run(); err != nil && err != nethttp.ErrServerClosed {
			config.ContextLogger.Fatal("server listening error:", err)
		}
	}()

	// catch shutdown
	errc := make(chan error, 1)

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		errc <- fmt.Errorf("%v", <-sigc)
	}()

	<-errc
	// First we close the connection with gorm:
	sqlDB, err := gormDb.DB()
	if err := sqlDB.Close(); err != nil {
		config.ContextLogger.WithError(err).Fatal("close gorm connection error")
	}
	// Second we close the connection with reids:
	if err := rcc.Close(); err != nil {
		config.ContextLogger.WithError(err).Fatal("close redis connection error")
	}

	// graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	srv.GracefulShutdown(ctx)
	<-ctx.Done()

	config.ContextLogger.Infoln("server exiting")
}
