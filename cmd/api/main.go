package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"go-api-server/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.Init(); err != nil {
		slog.Error("Failed to initialize configuration", "error", err)
		os.Exit(1)
	}

	cfg := config.Instance()
	if cfg == nil {
		slog.Error("Failed to get configuration instance")
		os.Exit(1)
	}

	gin.SetMode(cfg.Server.RunMode)
	r := gin.Default()
	r.GET("/hello-world", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	slog.Info("Starting server", "address", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
