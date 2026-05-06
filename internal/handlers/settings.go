package handlers

import (
	"pos-backend/internal/repository"
	"pos-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
	repo *repository.SettingsRepository
}

func NewSettingsHandler(repo *repository.SettingsRepository) *SettingsHandler {
	return &SettingsHandler{repo: repo}
}

func (h *SettingsHandler) GetSettings(c *gin.Context) {
	settings, err := h.repo.Find()
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	if settings == nil {
		// Return default settings
		utils.SuccessResponse(c, map[string]interface{}{
			"taxRate":          10,
			"serviceCharge":    0,
			"kitchenPrint":     true,
			"autoAcceptOrders": false,
			"soundEnabled":     true,
			"theme":            "light",
		})
		return
	}

	utils.SuccessResponse(c, settings)
}

func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	var settings map[string]interface{}
	if err := c.ShouldBindJSON(&settings); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	if err := h.repo.Update(settings); err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, settings)
}
