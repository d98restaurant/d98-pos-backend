package handlers

import (
	"pos-backend/internal/repository"
	"pos-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	userRepo *repository.UserRepository
}

func NewAdminHandler(userRepo *repository.UserRepository) *AdminHandler {
	return &AdminHandler{userRepo: userRepo}
}

// ChangeUserPassword - Admin can change any user's password
func (h *AdminHandler) ChangeUserPassword(c *gin.Context) {
	var req struct {
		UserID      string `json:"userId"`
		NewPassword string `json:"newPassword"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request")
		return
	}
	
	if req.NewPassword == "" || len(req.NewPassword) < 6 {
		utils.ErrorResponse(c, 400, "Password must be at least 6 characters")
		return
	}
	
	// Check if user exists
	user, err := h.userRepo.FindByID(req.UserID)
	if err != nil || user == nil {
		utils.ErrorResponse(c, 404, "User not found")
		return
	}
	
	// Update password directly
	err = h.userRepo.Update(req.UserID, map[string]interface{}{
		"passwordHash": req.NewPassword,
	})
	
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to update password")
		return
	}
	
	utils.SuccessResponse(c, gin.H{
		"message":  "Password changed successfully",
		"username": user.Username,
	})
}
