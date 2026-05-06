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
			"name":                "",
			"address":             "",
			"phone":               "",
			"email":               "",
			"gst":                 "",
			"fssai":               "",
			"upiId":               "",
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

	// Clean up the response - remove nested data structures
	cleanBusiness := map[string]interface{}{
		"name":                getString(business, "name"),
		"address":             getString(business, "address"),
		"phone":               getString(business, "phone"),
		"email":               getString(business, "email"),
		"gst":                 getString(business, "gst"),
		"fssai":               getString(business, "fssai"),
		"upiId":               getString(business, "upiId"),
		"currencySymbol":      getString(business, "currencySymbol"),
		"taxLabel":            getString(business, "taxLabel"),
		"footerMessage":       getString(business, "footerMessage"),
		"printBusinessName":   getBool(business, "printBusinessName", true),
		"printAddress":        getBool(business, "printAddress", true),
		"printPhone":          getBool(business, "printPhone", true),
		"printEmail":          getBool(business, "printEmail", true),
		"printGst":            getBool(business, "printGst", true),
		"printFssai":          getBool(business, "printFssai", true),
		"printHeaderDivider":  getBool(business, "printHeaderDivider", true),
		"printItems":          getBool(business, "printItems", true),
		"printTaxBreakdown":   getBool(business, "printTaxBreakdown", true),
		"printServiceCharge":  getBool(business, "printServiceCharge", true),
		"printGatewayCharges": getBool(business, "printGatewayCharges", true),
		"printFooter":         getBool(business, "printFooter", true),
		"printQrCode":         getBool(business, "printQrCode", true),
	}

	utils.SuccessResponse(c, cleanBusiness)
}

func (h *BusinessHandler) UpdateBusiness(c *gin.Context) {
	var business map[string]interface{}
	if err := c.ShouldBindJSON(&business); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body: "+err.Error())
		return
	}

	// Clean the data before saving
	cleanBusiness := map[string]interface{}{
		"name":                getString(business, "name"),
		"address":             getString(business, "address"),
		"phone":               getString(business, "phone"),
		"email":               getString(business, "email"),
		"gst":                 getString(business, "gst"),
		"fssai":               getString(business, "fssai"),
		"upiId":               getString(business, "upiId"),
		"currencySymbol":      getString(business, "currencySymbol"),
		"taxLabel":            getString(business, "taxLabel"),
		"footerMessage":       getString(business, "footerMessage"),
		"printBusinessName":   getBool(business, "printBusinessName", true),
		"printAddress":        getBool(business, "printAddress", true),
		"printPhone":          getBool(business, "printPhone", true),
		"printEmail":          getBool(business, "printEmail", true),
		"printGst":            getBool(business, "printGst", true),
		"printFssai":          getBool(business, "printFssai", true),
		"printHeaderDivider":  getBool(business, "printHeaderDivider", true),
		"printItems":          getBool(business, "printItems", true),
		"printTaxBreakdown":   getBool(business, "printTaxBreakdown", true),
		"printServiceCharge":  getBool(business, "printServiceCharge", true),
		"printGatewayCharges": getBool(business, "printGatewayCharges", true),
		"printFooter":         getBool(business, "printFooter", true),
		"printQrCode":         getBool(business, "printQrCode", true),
	}

	if err := h.repo.Update(cleanBusiness); err != nil {
		utils.ErrorResponse(c, 500, "Failed to save business details: "+err.Error())
		return
	}

	utils.SuccessResponse(c, cleanBusiness)
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getBool(m map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := m[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return defaultValue
}
