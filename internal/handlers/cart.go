package handlers

import (
	"pos-backend/internal/repository"
	"pos-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	repo *repository.CartRepository
}

func NewCartHandler(repo *repository.CartRepository) *CartHandler {
	return &CartHandler{repo: repo}
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.GetString("userID")
	cart, err := h.repo.FindByUserID(userID)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	if cart == nil {
		utils.SuccessResponse(c, gin.H{"items": []interface{}{}, "specialInstructions": map[string]string{}})
		return
	}

	utils.SuccessResponse(c, cart)
}

func (h *CartHandler) SaveCart(c *gin.Context) {
	userID := c.GetString("userID")
	var cartData map[string]interface{}
	if err := c.ShouldBindJSON(&cartData); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	if err := h.repo.Save(userID, cartData); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Cart saved successfully"})
}

func (h *CartHandler) AddItem(c *gin.Context) {
	userID := c.GetString("userID")
	var req struct {
		Item map[string]interface{} `json:"item"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	cart, err := h.repo.AddItem(userID, req.Item)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, cart)
}

func (h *CartHandler) UpdateItemQuantity(c *gin.Context) {
	userID := c.GetString("userID")
	itemID := c.Param("itemId")
	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	cart, err := h.repo.UpdateItemQuantity(userID, itemID, req.Quantity)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, cart)
}

func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID := c.GetString("userID")
	itemID := c.Param("itemId")

	cart, err := h.repo.RemoveItem(userID, itemID)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, cart)
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := c.GetString("userID")

	if err := h.repo.Clear(userID); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Cart cleared successfully"})
}
