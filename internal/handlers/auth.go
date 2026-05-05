package handlers

import (
	"pos-backend/internal/models"
	"pos-backend/internal/services"
	"pos-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	resp, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		utils.UnauthorizedResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, resp)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.CreatedResponse(c, resp)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	userID := c.GetString("userID")
	if err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Password changed successfully"})
}

func (h *AuthHandler) GetUsers(c *gin.Context) {
	users, err := h.authService.GetUsers()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, users)
}

func (h *AuthHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	if err := h.authService.UpdateUser(userID, &req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "User updated successfully"})
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	currentUserID := c.GetString("userID")
	if userID == currentUserID {
		utils.BadRequestResponse(c, "Cannot delete your own account")
		return
	}

	if err := h.authService.DeleteUser(userID); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "User deleted successfully"})
}
