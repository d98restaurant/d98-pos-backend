package handlers

import (
	"pos-backend/internal/repository"
	"pos-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type ResetHandler struct {
	userRepo *repository.UserRepository
}

func NewResetHandler(userRepo *repository.UserRepository) *ResetHandler {
	return &ResetHandler{userRepo: userRepo}
}

func (h *ResetHandler) ClearAndReset(c *gin.Context) {
	secret := c.Query("secret")
	
	if secret != "reset2024" {
		utils.ErrorResponse(c, 403, "Unauthorized")
		return
	}
	
	// Clear all user data
	err := repository.ClearAllUsers()
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	
	// Create a default admin user
	defaultUser := map[string]interface{}{
		"username": "admin",
		"email":    "admin@pos.com",
		"password": "admin123",
		"role":     "admin",
	}
	
	utils.SuccessResponse(c, defaultUser)
}
