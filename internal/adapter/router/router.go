package router

import (
	"net/http"
	"time"

	"go-api-server/internal/adapter/handler"
	"go-api-server/internal/app"
	"go-api-server/internal/pkg/timeutil"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(engine *gin.Engine, core *app.CoreServices) {
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": timeutil.Now().Format(time.RFC3339Nano),
		})
	})

	v1 := engine.Group("/v1")
	{
		setupConfigurationRoutes(v1, core)
	}
}

func setupConfigurationRoutes(r *gin.RouterGroup, core *app.CoreServices) {
	configurationHandler := handler.NewConfigurationHandler(core.AppServices.ConfigurationService)

	configRoutes := r.Group("/configurations")
	{
		configRoutes.GET("/app-version", configurationHandler.GetApplicationVersion)
	}
}
