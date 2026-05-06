package handlers

import (
	"pos-backend/internal/repository"
	"pos-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type BusinessHandler struct {
	repo *repository.BusinessRepository
}

func NewBusinessHandler(repo *repository.BusinessRepository) *BusinessHandler {
	return &BusinessHandler{repo: repo}
}

func (h *BusinessHandler) GetBusiness(c *gin.Context) {
	business, err := h.repo.Find()
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	if business == nil {
		// Return default business details
		utils.SuccessResponse(c, map[string]interface{}{
			"name":                "My Restaurant",
			"address":             "123 Main Street",
			"phone":               "+91 9876543210",
			"email":               "info@restaurant.com",
			"gst":                 "27ABCDE1234F1Z5",
			"fssai":               "12345678901234",
			"upiId":               "paytm.s1yxcay@pty",
			"currencySymbol":      "₹",
			"taxLabel":            "GST",
			"footerMessage":       "Thank you! Visit Again!",
			"printBusinessName":   true,
			"printAddress":        true,
			"printPhone":          true,
			"printEmail":          true,
			"printGst":            true,
			"printFssai":          true,
			"printHeaderDivider":  true,
			"printItems":          true,
			"printTaxBreakdown":   true,
			"printServiceCharge":  true,
			"printGatewayCharges": true,
			"printFooter":         true,
			"printQrCode":         true,
		})
		return
	}

	utils.SuccessResponse(c, business)
}

func (h *BusinessHandler) UpdateBusiness(c *gin.Context) {
	var business map[string]interface{}
	if err := c.ShouldBindJSON(&business); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body: "+err.Error())
		return
	}

	if err := h.repo.Update(business); err != nil {
		utils.ErrorResponse(c, 500, "Failed to save business details: "+err.Error())
		return
	}

	utils.SuccessResponse(c, business)
}
