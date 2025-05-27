// @title Go API ʕ ´•ᴥ•`ʔﾉﾞ
// @version 1.0
// @description This is a Go API server.
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "go-api-server/docs"
	"go-api-server/internal/adapter/router"
	"go-api-server/internal/app"
	"go-api-server/internal/config"
	"go-api-server/internal/pkg/logger"
	"go-api-server/internal/pkg/timeouts"

	"github.com/gin-gonic/gin"
)

var log = logger.Instance()

func main() {
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	defer logger.Shutdown()

	if err := config.Init(); err != nil {
		log.WithError(err).Fatal("Configuration initialization failed")
	}
	cfg := config.Instance()

	logger.Configure(cfg.Server.Logger)

	log.Info("Application bootstrap starting...")

	coreServices, coreCleanup, err := app.NewCoreServices(rootCtx, cfg)
	if err != nil {
		log.WithError(err).Error("Failed to initialize core services")
		rootCancel()
		if coreCleanup != nil {
			coreCleanup()
		}
		log.Info("Application shutting down due to core services initialization failure")
		os.Exit(1)
	}
	if coreCleanup != nil {
		defer coreCleanup()
	}

	gin.SetMode(cfg.Server.RunMode)
	engine := gin.Default()
	router.SetupRoutes(engine, coreServices)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      engine,
		ReadTimeout:  timeouts.ServerRead(cfg),
		WriteTimeout: timeouts.ServerWrite(cfg),
		IdleTimeout:  timeouts.ServerIdle(cfg),
	}

	serverErrChan := make(chan error, 1)
	go func() {
		log.Infof("HTTP server starting, listening on port %d", cfg.Server.Port)
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			serverErrChan <- fmt.Errorf("HTTP server ListenAndServe failed: %w", err)
		}
	}()

	quitSignalChan := make(chan os.Signal, 1)
	signal.Notify(quitSignalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quitSignalChan:
		log.Infof("Received OS signal: %s, initiating graceful shutdown", sig)
	case err := <-serverErrChan:
		log.WithError(err).Error("HTTP server error detected, initiating shutdown")
	case <-rootCtx.Done():
		log.Info("Root context cancelled, initiating graceful shutdown")
	}

	signal.Stop(quitSignalChan)
	rootCancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), timeouts.ServerShutdown(cfg))
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.WithError(err).Error("Graceful HTTP server shutdown failed or timed out")
	} else {
		log.Info("HTTP server gracefully stopped")
	}

	log.Info("Application shutdown complete")
}
