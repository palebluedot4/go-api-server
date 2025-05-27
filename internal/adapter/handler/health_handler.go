package handler

import (
	"net/http"
	"time"

	"go-api-server/internal/pkg/timeutil"

	"github.com/gin-gonic/gin"
)

type HealthCheckHandler struct{}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

// HealthCheck godoc
// @Summary Health check
// @Description Checks the health of the application.
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string "Successfully returned health status"
// @Router /health [get]
func (h *HealthCheckHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": timeutil.Now().Format(time.RFC3339Nano),
	})
}
