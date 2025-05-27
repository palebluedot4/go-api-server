package router

import (
	"go-api-server/internal/adapter/handler"
	"go-api-server/internal/app"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(engine *gin.Engine, core *app.CoreServices) {
	engine.GET("/health", handler.NewHealthCheckHandler().HealthCheck)

	v1 := engine.Group("/v1")
	{
		setupConfigurationRoutes(v1, core)
	}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func setupConfigurationRoutes(r *gin.RouterGroup, core *app.CoreServices) {
	configurationHandler := handler.NewConfigurationHandler(core.AppServices.ConfigurationService)

	configRoutes := r.Group("/configurations")
	{
		configRoutes.GET("/app-version", configurationHandler.GetApplicationVersion)
	}
}
