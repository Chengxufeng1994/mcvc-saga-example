package main

import (
	"context"
	"fmt"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Chengxufeng1994/go-saga-example/product-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/di"

	log "github.com/sirupsen/logrus"
)

func main() {
	bootCfg := di.InitBootstrapConfig("config/orchestrator")
	appCfg := di.InitApplicationConfig("config/orchestrator")
	config.InitLogger(appCfg, bootCfg)

	srv := di.InitializeOrchestratorServer(appCfg, bootCfg)

	go func() {
		if err := srv.Run(); err != nil && err != nethttp.ErrServerClosed {
			log.Fatal("server listening error:", err)
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

	// graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	srv.GracefulShutdown(ctx)
	<-ctx.Done()

	log.Infoln("server exiting")
}
