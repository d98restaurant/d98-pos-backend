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
	orders, err := h.orderService.GetAllOrders(c.Request.Context())
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}
	utils.SuccessResponse(c, orders)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	order, err := h.orderService.GetOrderByID(c.Request.Context(), id)
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

	createdOrder, err := h.orderService.CreateOrder(c.Request.Context(), &order)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	// Broadcast to WebSocket clients
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

	updatedOrder, err := h.orderService.UpdateOrderStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	// Broadcast update
	h.wsHub.BroadcastToRoom("orders", "order-updated", updatedOrder)
	if req.Status == models.OrderStatusAccepted {
		h.wsHub.BroadcastToRoom("orders", "order-accepted", id)
	} else if req.Status == models.OrderStatusReadyForBilling {
		h.wsHub.BroadcastToRoom("orders", "order-ready-for-billing", id)
	} else if req.Status == models.OrderStatusCompleted {
		h.wsHub.BroadcastToRoom("orders", "order-completed", id)
	}

	utils.SuccessResponse(c, updatedOrder)
}

func (h *OrderHandler) CompletePayment(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		PaymentMethod   string                 `json:"paymentMethod"`
		PaymentDetails  map[string]interface{} `json:"paymentDetails"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	updatedOrder, err := h.orderService.CompletePayment(c.Request.Context(), id, req.PaymentMethod, req.PaymentDetails)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	h.wsHub.BroadcastToRoom("orders", "order-completed", id)
	h.wsHub.BroadcastToRoom("orders", "order-updated", updatedOrder)

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

	updatedOrder, err := h.orderService.AddItemToOrder(c.Request.Context(), id, &req.Item)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	h.wsHub.BroadcastToRoom("orders", "order-item-added", map[string]interface{}{
		"orderId": id,
		"item":    req.Item,
	})
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

	updatedOrder, err := h.orderService.UpdateItemQuantity(c.Request.Context(), orderID, itemID, req.Quantity)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	h.wsHub.BroadcastToRoom("orders", "order-item-quantity-updated", map[string]interface{}{
		"orderId":     orderID,
		"itemId":      itemID,
		"newQuantity": req.Quantity,
	})
	h.wsHub.BroadcastToRoom("orders", "order-updated", updatedOrder)

	utils.SuccessResponse(c, updatedOrder)
}

func (h *OrderHandler) RemoveItemFromOrder(c *gin.Context) {
	orderID := c.Param("id")
	itemID := c.Param("itemId")

	updatedOrder, err := h.orderService.RemoveItemFromOrder(c.Request.Context(), orderID, itemID)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	h.wsHub.BroadcastToRoom("orders", "order-item-removed", map[string]interface{}{
		"orderId": orderID,
		"itemId":  itemID,
	})
	h.wsHub.BroadcastToRoom("orders", "order-updated", updatedOrder)

	utils.SuccessResponse(c, updatedOrder)
}

func (h *OrderHandler) UpdateItemStatus(c *gin.Context) {
	orderID := c.Param("id")
	itemID := c.Param("itemId")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	updatedOrder, err := h.orderService.UpdateItemStatus(c.Request.Context(), orderID, itemID, req.Status)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	h.wsHub.BroadcastToRoom("orders", "item-status-updated", updatedOrder)

	utils.SuccessResponse(c, updatedOrder)
}

func (h *OrderHandler) GetActiveOrdersByTable(c *gin.Context) {
	tableNumber, err := strconv.Atoi(c.Param("tableNumber"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid table number")
		return
	}

	orders, err := h.orderService.GetActiveOrdersByTable(c.Request.Context(), tableNumber)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, orders)
}

func (h *OrderHandler) CompleteTableBilling(c *gin.Context) {
	tableNumber, err := strconv.Atoi(c.Param("tableNumber"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid table number")
		return
	}

	var req struct {
		SessionID string `json:"sessionId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	if err := h.orderService.CompleteTableBilling(c.Request.Context(), tableNumber, req.SessionID); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Table billing completed"})
}

func (h *OrderHandler) GetPendingCancellationRequests(c *gin.Context) {
	requests, err := h.orderService.GetPendingCancellationRequests(c.Request.Context())
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}
	utils.SuccessResponse(c, requests)
}

func (h *OrderHandler) RequestItemCancellation(c *gin.Context) {
	orderID := c.Param("id")
	itemID := c.Param("itemId")
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Reason = ""
	}

	err := h.orderService.RequestItemCancellation(c.Request.Context(), orderID, itemID, req.Reason)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	h.wsHub.BroadcastToRoom("orders", "cancellation-requested", map[string]interface{}{
		"orderId": orderID,
		"itemId":  itemID,
		"reason":  req.Reason,
	})

	utils.SuccessResponse(c, gin.H{"message": "Cancellation requested"})
}

func (h *OrderHandler) ApproveCancellation(c *gin.Context) {
	orderID := c.Param("id")
	itemID := c.Param("itemId")

	updatedOrder, err := h.orderService.ApproveCancellation(c.Request.Context(), orderID, itemID)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	h.wsHub.BroadcastToRoom("orders", "cancellation-approved", map[string]interface{}{
		"orderId": orderID,
		"itemId":  itemID,
	})

	utils.SuccessResponse(c, updatedOrder)
}

func (h *OrderHandler) RejectCancellation(c *gin.Context) {
	orderID := c.Param("id")
	itemID := c.Param("itemId")
	var req struct {
		RejectReason string `json:"rejectReason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		req.RejectReason = ""
	}

	err := h.orderService.RejectCancellation(c.Request.Context(), orderID, itemID, req.RejectReason)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	h.wsHub.BroadcastToRoom("orders", "cancellation-rejected", map[string]interface{}{
		"orderId": orderID,
		"itemId":  itemID,
		"reason":  req.RejectReason,
	})

	utils.SuccessResponse(c, gin.H{"message": "Cancellation rejected"})
}

func (h *OrderHandler) GetCreditCustomers(c *gin.Context) {
	customers, err := h.orderService.GetCreditCustomers(c.Request.Context())
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}
	utils.SuccessResponse(c, customers)
}

func (h *OrderHandler) ProcessCreditCollection(c *gin.Context) {
	var req struct {
		CustomerID    string  `json:"customerId"`
		Amount        float64 `json:"amount"`
		PaymentMethod string  `json:"paymentMethod"`
		Note          string  `json:"note"`
		CollectedBy   string  `json:"collectedBy"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	if err := h.orderService.ProcessCreditCollection(c.Request.Context(), req.CustomerID, req.Amount, req.PaymentMethod, req.Note, req.CollectedBy); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Payment collected successfully"})
}

func (h *OrderHandler) ChangePaymentMethod(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		PaymentMethod string `json:"paymentMethod"`
		Reason        string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	updatedOrder, err := h.orderService.ChangePaymentMethod(c.Request.Context(), id, req.PaymentMethod, req.Reason)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	h.wsHub.BroadcastToRoom("orders", "order-updated", updatedOrder)

	utils.SuccessResponse(c, updatedOrder)
}