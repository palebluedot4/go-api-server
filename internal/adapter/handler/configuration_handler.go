package handler

import (
	"net/http"

	"go-api-server/internal/adapter/handler/dto/response"
	"go-api-server/internal/domain/service"

	"github.com/gin-gonic/gin"
)

type ConfigurationHandler struct {
	configurationService service.ConfigurationService
}

func NewConfigurationHandler(configurationService service.ConfigurationService) *ConfigurationHandler {
	return &ConfigurationHandler{
		configurationService: configurationService,
	}
}

// GetApplicationVersion handles the HTTP GET request to retrieve the application version.
// @Summary Get application version
// @Description Retrieves the application version.
// @Tags configurations
// @Produce json
// @Success 200 {object} response.GetApplicationVersionResponse "Successfully retrieved application version"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/configurations/app-version [get]
func (h *ConfigurationHandler) GetApplicationVersion(c *gin.Context) {
	const appVersionKey = "app_version"

	entity, err := h.configurationService.GetConfigurationByKey(c, appVersionKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp := response.GetApplicationVersionResponse{
		Version: entity.ConfigValue,
	}
	c.JSON(http.StatusOK, resp)
}
