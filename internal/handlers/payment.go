package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"pos-backend/internal/services"
	"pos-backend/internal/utils"

	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
	orderService   *services.OrderService
}

func NewPaymentHandler(paymentService *services.PaymentService, orderService *services.OrderService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		orderService:   orderService,
	}
}

func (h *PaymentHandler) CreateRazorpayOrder(c *gin.Context) {
	var req struct {
		Amount  int    `json:"amount"`
		Currency string `json:"currency"`
		Receipt  string `json:"receipt"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	client := razorpay.NewClient(h.paymentService.GetKeyID(), h.paymentService.GetKeySecret())
	
	data := map[string]interface{}{
		"amount":   req.Amount,
		"currency": req.Currency,
		"receipt":  req.Receipt,
	}
	
	body, err := client.Order.Create(data, nil)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, body)
}

func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	var req struct {
		OrderID   string `json:"orderId"`
		PaymentID string `json:"paymentId"`
		Signature string `json:"signature"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	secret := h.paymentService.GetKeySecret()
	message := req.OrderID + "|" + req.PaymentID
	hx := hmac.New(sha256.New, []byte(secret))
	hx.Write([]byte(message))
	expectedSignature := hex.EncodeToString(hx.Sum(nil))

	if expectedSignature != req.Signature {
		utils.BadRequestResponse(c, "Invalid payment signature")
		return
	}

	utils.SuccessResponse(c, gin.H{"success": true})
}

func (h *PaymentHandler) ProcessCreditSale(c *gin.Context) {
	var req struct {
		OrderID      string  `json:"orderId"`
		CustomerName string  `json:"customerName"`
		CustomerPhone string `json:"customerPhone"`
		CustomerEmail string `json:"customerEmail"`
		DueDate      *time.Time `json:"dueDate"`
		Amount       float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	paymentDetails := map[string]interface{}{
		"method":       "credit",
		"amount":       req.Amount,
		"status":       "credit_due",
		"customerName": req.CustomerName,
		"customerPhone": req.CustomerPhone,
		"customerEmail": req.CustomerEmail,
		"transactionId": fmt.Sprintf("CREDIT_%d_%d", time.Now().UnixNano(), rand.Intn(10000)),
	}
	if req.DueDate != nil {
		paymentDetails["dueDate"] = req.DueDate
	}

	updatedOrder, err := h.orderService.CompletePayment(c.Request.Context(), req.OrderID, "credit", paymentDetails)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, updatedOrder)
}