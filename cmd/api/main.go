package main

import (
	"fmt"
	"net/http"

	"go-api-server/internal/config"
	"go-api-server/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

var log = logger.Instance()

func main() {
	if err := config.Init(); err != nil {
		log.WithError(err).Fatal("Failed to initialize configuration")
	}
	cfg := config.Instance()

	logger.Configure(cfg.Server.Logger)
	defer logger.Shutdown()

	log.WithField("mode", cfg.Server.RunMode).Info("Gin run mode set")

	gin.SetMode(cfg.Server.RunMode)
	r := gin.Default()
	r.GET("/hello-world", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.WithField("address", serverAddr).Info("Starting server")
	if err := r.Run(serverAddr); err != nil {
		log.WithError(err).Fatal("Failed to start server")
	}
}
