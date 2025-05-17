package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-api-server/internal/config"
	"go-api-server/internal/pkg/logger"
	"go-api-server/internal/pkg/timeouts"
	"go-api-server/internal/pkg/timeutil"

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

	log.WithFields(map[string]any{
		"port":      cfg.Server.Port,
		"run_mode":  cfg.Server.RunMode,
		"log_level": cfg.Server.Logger.Level,
	}).Info("Application starting with loaded configuration")

	gin.SetMode(cfg.Server.RunMode)
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": timeutil.Now().Format(time.RFC3339Nano),
		})
	})

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  timeouts.ServerRead(cfg),
		WriteTimeout: timeouts.ServerWrite(cfg),
		IdleTimeout:  timeouts.ServerIdle(cfg),
	}

	serverErrChan := make(chan error, 1)
	go func() {
		err := s.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
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

	rootCancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), timeouts.ServerShutdown(cfg))
	defer shutdownCancel()

	log.Info("Attempting graceful HTTP server shutdown")
	if err := s.Shutdown(shutdownCtx); err != nil {
		log.WithError(err).Error("Graceful HTTP server shutdown failed or timed out")
	} else {
		log.Info("HTTP server gracefully stopped")
	}

	log.Info("Application shutdown complete")
}
