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
	settings, err := h.repo.Find(c.Request.Context())
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
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
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	if err := h.repo.Update(c.Request.Context(), settings); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, settings)
}