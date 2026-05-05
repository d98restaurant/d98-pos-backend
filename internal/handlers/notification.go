package handlers

import (
	"pos-backend/internal/services"
	"pos-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service *services.NotificationService
}

func NewNotificationHandler(service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) Subscribe(c *gin.Context) {
	var subscription map[string]interface{}
	if err := c.ShouldBindJSON(&subscription); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	userID := c.GetString("userID")
	role := c.GetString("role")

	if err := h.service.Subscribe(c.Request.Context(), userID, role, subscription); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Subscribed successfully", "kitchenSubscribers": h.service.GetKitchenSubscriberCount()})
}

func (h *NotificationHandler) Unsubscribe(c *gin.Context) {
	userID := c.GetString("userID")
	
	if err := h.service.Unsubscribe(c.Request.Context(), userID); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Unsubscribed successfully"})
}

func (h *NotificationHandler) GetVAPIDPublicKey(c *gin.Context) {
	utils.SuccessResponse(c, gin.H{"publicKey": h.service.GetVAPIDPublicKey()})
}
