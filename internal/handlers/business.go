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
	business, err := h.repo.Find(c.Request.Context())
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	if business == nil {
		// Return default business details
		utils.SuccessResponse(c, map[string]interface{}{
			"name":                "Restaurant POS",
			"address":             "",
			"phone":               "",
			"email":               "",
			"gst":                 "",
			"fssai":               "",
			"upiId":               "paytm.s1yxcay@pty",
			"currencySymbol":      "₹",
			"taxLabel":            "GST",
			"footerMessage":       "Thank you! Visit Again!",
			"printBusinessName":   true,
			"printAddress":        true,
			"printPhone":          true,
			"printGst":            true,
			"printItems":          true,
			"printTaxBreakdown":   true,
			"printServiceCharge":  true,
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
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	if err := h.repo.Update(c.Request.Context(), business); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, business)
}