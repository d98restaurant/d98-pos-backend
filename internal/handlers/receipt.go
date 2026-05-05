package handlers

import (
	"pos-backend/internal/repository"
	"pos-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type ReceiptHandler struct {
	receiptRepo  *repository.ReceiptRepository
	businessRepo *repository.BusinessRepository
}

func NewReceiptHandler(receiptRepo *repository.ReceiptRepository, businessRepo *repository.BusinessRepository) *ReceiptHandler {
	return &ReceiptHandler{
		receiptRepo:  receiptRepo,
		businessRepo: businessRepo,
	}
}

func (h *ReceiptHandler) GetReceipt(c *gin.Context) {
	id := c.Param("id")
	receipt, err := h.receiptRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}
	if receipt == nil {
		utils.NotFoundResponse(c, "Receipt not found")
		return
	}

	utils.SuccessResponse(c, receipt)
}

func (h *ReceiptHandler) CreateReceipt(c *gin.Context) {
	var receipt map[string]interface{}
	if err := c.ShouldBindJSON(&receipt); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	if err := h.receiptRepo.Create(c.Request.Context(), receipt); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.CreatedResponse(c, receipt)
}