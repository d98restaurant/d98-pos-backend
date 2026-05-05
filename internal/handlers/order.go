package handlers

import (
	"strconv"

	"pos-backend/internal/models"
	"pos-backend/internal/services"
	"pos-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *services.OrderService
	wsHub        *utils.Hub
}

func NewOrderHandler(orderService *services.OrderService, wsHub *utils.Hub) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		wsHub:        wsHub,
	}
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	orders, err := h.orderService.GetAllOrders()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}
	utils.SuccessResponse(c, orders)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	order, err := h.orderService.GetOrderByID(id)
	if err != nil {
		utils.NotFoundResponse(c, err.Error())
		return
	}
	utils.SuccessResponse(c, order)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	createdOrder, err := h.orderService.CreateOrder(&order)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	h.wsHub.BroadcastToRoom("orders", "new-order-received", createdOrder)
	utils.CreatedResponse(c, createdOrder)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status models.OrderStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	updatedOrder, err := h.orderService.UpdateOrderStatus(id, req.Status)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	h.wsHub.BroadcastToRoom("orders", "order-updated", updatedOrder)
	utils.SuccessResponse(c, updatedOrder)
}

func (h *OrderHandler) CompletePayment(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		PaymentMethod  string                 `json:"paymentMethod"`
		PaymentDetails map[string]interface{} `json:"paymentDetails"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	updatedOrder, err := h.orderService.CompletePayment(id, req.PaymentMethod, req.PaymentDetails)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	h.wsHub.BroadcastToRoom("orders", "order-completed", id)
	utils.SuccessResponse(c, updatedOrder)
}

func (h *OrderHandler) AddItemToOrder(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Item models.OrderItem `json:"item"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	updatedOrder, err := h.orderService.AddItemToOrder(id, &req.Item)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	h.wsHub.BroadcastToRoom("orders", "order-updated", updatedOrder)
	utils.SuccessResponse(c, updatedOrder)
}

func (h *OrderHandler) UpdateItemQuantity(c *gin.Context) {
	orderID := c.Param("id")
	itemID := c.Param("itemId")
	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	updatedOrder, err := h.orderService.UpdateItemQuantity(orderID, itemID, req.Quantity)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	h.wsHub.BroadcastToRoom("orders", "order-updated", updatedOrder)
	utils.SuccessResponse(c, updatedOrder)
}

func (h *OrderHandler) RemoveItemFromOrder(c *gin.Context) {
	orderID := c.Param("id")
	itemID := c.Param("itemId")

	updatedOrder, err := h.orderService.RemoveItemFromOrder(orderID, itemID)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	h.wsHub.BroadcastToRoom("orders", "order-updated", updatedOrder)
	utils.SuccessResponse(c, updatedOrder)
}

// Placeholder methods for other endpoints
func (h *OrderHandler) UpdateItemStatus(c *gin.Context) { utils.SuccessResponse(c, gin.H{}) }
func (h *OrderHandler) GetActiveOrdersByTable(c *gin.Context) { utils.SuccessResponse(c, []interface{}{}) }
func (h *OrderHandler) CompleteTableBilling(c *gin.Context) { utils.SuccessResponse(c, gin.H{}) }
func (h *OrderHandler) GetPendingCancellationRequests(c *gin.Context) { utils.SuccessResponse(c, []interface{}{}) }
func (h *OrderHandler) RequestItemCancellation(c *gin.Context) { utils.SuccessResponse(c, gin.H{}) }
func (h *OrderHandler) ApproveCancellation(c *gin.Context) { utils.SuccessResponse(c, gin.H{}) }
func (h *OrderHandler) RejectCancellation(c *gin.Context) { utils.SuccessResponse(c, gin.H{}) }
func (h *OrderHandler) GetCreditCustomers(c *gin.Context) { utils.SuccessResponse(c, []interface{}{}) }
func (h *OrderHandler) ProcessCreditCollection(c *gin.Context) { utils.SuccessResponse(c, gin.H{}) }
func (h *OrderHandler) ChangePaymentMethod(c *gin.Context) { utils.SuccessResponse(c, gin.H{}) }
