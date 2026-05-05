package handlers

import (
	"pos-backend/internal/repository"
	"pos-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type DebugHandler struct {
	userRepo *repository.UserRepository
}

func NewDebugHandler(userRepo *repository.UserRepository) *DebugHandler {
	return &DebugHandler{userRepo: userRepo}
}

func (h *DebugHandler) CheckUser(c *gin.Context) {
	username := c.Query("username")
	secret := c.Query("secret")
	
	// Security: Only allow with secret
	if secret != "debug2024" {
		utils.ErrorResponse(c, 403, "Unauthorized")
		return
	}
	
	user, err := h.userRepo.FindByUsername(username)
	if err != nil || user == nil {
		utils.ErrorResponse(c, 404, "User not found")
		return
	}
	
	// Test the password "mypassword" against stored hash
	testPassword := "mypassword"
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(testPassword))
	
	utils.SuccessResponse(c, gin.H{
		"username":      user.Username,
		"hash_length":   len(user.PasswordHash),
		"hash_prefix":   user.PasswordHash[:min(30, len(user.PasswordHash))],
		"is_bcrypt":     len(user.PasswordHash) > 3 && (user.PasswordHash[:3] == "$2a" || user.PasswordHash[:3] == "$2b"),
		"password_test": err == nil,
		"error":         err.Error(),
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
