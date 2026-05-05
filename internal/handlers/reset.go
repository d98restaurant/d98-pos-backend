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

func (h *ResetHandler) ForceResetPassword(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Secret   string `json:"secret"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request")
		return
	}

	// Security: Use a secret key to prevent unauthorized resets
	if req.Secret != "pos-reset-2024" {
		utils.ErrorResponse(c, 403, "Unauthorized")
		return
	}

	// Find user
	user, err := h.userRepo.FindByUsername(req.Username)
	if err != nil || user == nil {
		utils.ErrorResponse(c, 404, "User not found")
		return
	}

	// Hash new password
	newHash, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to hash password")
		return
	}

	// Update password
	err = h.userRepo.Update(user.ID, map[string]interface{}{
		"passwordHash": newHash,
	})

	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to update password")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message":  "Password reset successfully",
		"username": req.Username,
	})
}
