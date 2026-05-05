package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"pos-backend/internal/models"
	"pos-backend/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderService struct {
	orderRepo    *repository.OrderRepository
	tableRepo    *repository.TableRepository
	menuRepo     *repository.MenuRepository
	notification *NotificationService
}

func NewOrderService(orderRepo *repository.OrderRepository, tableRepo *repository.TableRepository, menuRepo *repository.MenuRepository, notification *NotificationService) *OrderService {
	return &OrderService{
		orderRepo:    orderRepo,
		tableRepo:    tableRepo,
		menuRepo:     menuRepo,
		notification: notification,
	}
}

func (s *OrderService) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	return s.orderRepo.FindAll(ctx)
}

func (s *OrderService) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	return s.orderRepo.FindByID(ctx, id)
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	// Generate order number
	orderNumber, err := s.orderRepo.GetNextOrderNumber(ctx)
	if err != nil {
		orderNumber = int(time.Now().UnixNano() % 1000000)
	}
	order.OrderNumber = orderNumber

	// Handle additional orders (running orders on same table)
	if order.IsAdditionalOrder && order.TableNumber > 0 {
		table, err := s.tableRepo.FindByNumber(ctx, order.TableNumber)
		if err == nil && table != nil && table.RunningOrderCount > 0 {
			order.RunningNumber = table.RunningOrderCount + 1
			order.BaseOrderNumber = order.OrderNumber
			order.DisplayOrderNumber = fmt.Sprintf("%d-%d", order.OrderNumber, order.RunningNumber)
		}
	}

	// Set initial status
	order.Status = models.OrderStatusPending

	// Create order
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	// Update table running count for dine-in
	if order.OrderType == models.OrderTypeDineIn && order.TableNumber > 0 {
		if err := s.tableRepo.IncrementRunningCount(ctx, order.TableNumber, order.Total); err != nil {
			// Log error but don't fail order creation
		}
	}

	// Send notification to kitchen
	s.notification.SendOrderNotification(order)

	return order, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, id string, status models.OrderStatus) (*models.Order, error) {
	if err := s.orderRepo.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}
	return s.orderRepo.FindByID(ctx, id)
}

func (s *OrderService) CompletePayment(ctx context.Context, id, paymentMethod string, paymentDetails map[string]interface{}) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	now := time.Now()
	paymentDetails["paidAt"] = now
	paymentDetails["method"] = paymentMethod
	paymentDetails["status"] = "completed"

	if err := s.orderRepo.Update(ctx, id, bson.M{
		"payment":   paymentDetails,
		"status":    models.OrderStatusCompleted,
		"completedAt": now,
	}); err != nil {
		return nil, err
	}

	// Update table running counts for dine-in
	if order.OrderType == models.OrderTypeDineIn && order.TableNumber > 0 {
		s.tableRepo.DecrementRunningCount(ctx, order.TableNumber, order.Total)
	}

	return s.orderRepo.FindByID(ctx, id)
}

func (s *OrderService) AddItemToOrder(ctx context.Context, orderID string, item *models.OrderItem) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	// Check if item already exists
	found := false
	for i, existingItem := range order.Items {
		if existingItem.ID == item.ID {
			// Update quantity and mark as modified
			order.Items[i].Quantity += item.Quantity
			order.Items[i].IsModified = true
			now := time.Now()
			order.Items[i].ModifiedAt = &now
			order.Items[i].OldQuantity = existingItem.Quantity
			found = true
			break
		}
	}

	if !found {
		item.IsModified = true
		now := time.Now()
		item.ModifiedAt = &now
		order.Items = append(order.Items, *item)
	}

	// Recalculate totals
	s.recalculateOrderTotals(order)

	order.HasModifications = true
	order.UpdatedAt = time.Now()

	if err := s.orderRepo.Update(ctx, orderID, bson.M{
		"items":             order.Items,
		"subtotal":          order.Subtotal,
		"tax":               order.Tax,
		"serviceCharge":     order.ServiceCharge,
		"total":             order.Total,
		"hasModifications":  true,
		"updatedAt":         order.UpdatedAt,
	}); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) UpdateItemQuantity(ctx context.Context, orderID, itemID string, quantity int) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	for i, item := range order.Items {
		if item.ID == itemID {
			if quantity <= 0 {
				// Remove item
				order.Items = append(order.Items[:i], order.Items[i+1:]...)
			} else {
				order.Items[i].Quantity = quantity
				order.Items[i].IsModified = true
				now := time.Now()
				order.Items[i].ModifiedAt = &now
			}
			break
		}
	}

	s.recalculateOrderTotals(order)
	order.UpdatedAt = time.Now()

	if err := s.orderRepo.Update(ctx, orderID, bson.M{
		"items":         order.Items,
		"subtotal":      order.Subtotal,
		"tax":           order.Tax,
		"serviceCharge": order.ServiceCharge,
		"total":         order.Total,
		"updatedAt":     order.UpdatedAt,
	}); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) RemoveItemFromOrder(ctx context.Context, orderID, itemID string) (*models.Order, error) {
	return s.UpdateItemQuantity(ctx, orderID, itemID, 0)
}

func (s *OrderService) UpdateItemStatus(ctx context.Context, orderID, itemID, status string) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	for i, item := range order.Items {
		if item.ID == itemID {
			order.Items[i].Status = status
			break
		}
	}

	// Check if all items are completed
	allCompleted := true
	for _, item := range order.Items {
		if item.Status != "completed" && !item.IsRemoved {
			allCompleted = false
			break
		}
	}

	if allCompleted && order.Status != models.OrderStatusCompleted {
		order.Status = models.OrderStatusReadyForBilling
	}

	order.UpdatedAt = time.Now()

	if err := s.orderRepo.Update(ctx, orderID, bson.M{
		"items":     order.Items,
		"status":    order.Status,
		"updatedAt": order.UpdatedAt,
	}); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetActiveOrdersByTable(ctx context.Context, tableNumber int) ([]models.Order, error) {
	return s.orderRepo.FindByTable(ctx, tableNumber, true)
}

func (s *OrderService) CompleteTableBilling(ctx context.Context, tableNumber int, sessionID string) error {
	orders, err := s.orderRepo.FindByTable(ctx, tableNumber, true)
	if err != nil {
		return err
	}

	for _, order := range orders {
		if order.TableSessionID == sessionID {
			if err := s.orderRepo.UpdateStatus(ctx, order.ID.Hex(), models.OrderStatusCompleted); err != nil {
				return err
			}
		}
	}

	return s.tableRepo.ResetRunningOrders(ctx, tableNumber)
}

func (s *OrderService) GetPendingCancellationRequests(ctx context.Context) ([]models.CancellationRequest, error) {
	// This would typically be stored in a separate collection
	// For now, return empty slice
	return []models.CancellationRequest{}, nil
}

func (s *OrderService) RequestItemCancellation(ctx context.Context, orderID, itemID, reason string) error {
	// Store cancellation request
	return nil
}

func (s *OrderService) ApproveCancellation(ctx context.Context, orderID, itemID string) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	for i, item := range order.Items {
		if item.ID == itemID {
			order.Items[i].IsRemoved = true
			order.Items[i].Status = "cancelled"
			break
		}
	}

	s.recalculateOrderTotals(order)
	order.UpdatedAt = time.Now()

	if err := s.orderRepo.Update(ctx, orderID, bson.M{
		"items":     order.Items,
		"subtotal":  order.Subtotal,
		"tax":       order.Tax,
		"total":     order.Total,
		"updatedAt": order.UpdatedAt,
	}); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) RejectCancellation(ctx context.Context, orderID, itemID, reason string) error {
	// Just log the rejection, no changes to order
	return nil
}

func (s *OrderService) GetCreditCustomers(ctx context.Context) ([]map[string]interface{}, error) {
	return s.orderRepo.FindCreditCustomers(ctx)
}

func (s *OrderService) ProcessCreditCollection(ctx context.Context, customerID string, amount float64, paymentMethod, note, collectedBy string) error {
	// Update the credit collection record
	return nil
}

func (s *OrderService) ChangePaymentMethod(ctx context.Context, orderID, paymentMethod, reason string) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	if order.Payment == nil {
		order.Payment = &models.PaymentDetails{}
	}
	order.Payment.Method = paymentMethod
	order.UpdatedAt = time.Now()

	if err := s.orderRepo.Update(ctx, orderID, bson.M{
		"payment":   order.Payment,
		"updatedAt": order.UpdatedAt,
	}); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) recalculateOrderTotals(order *models.Order) {
	subtotal := 0.0
	for _, item := range order.Items {
		if !item.IsRemoved {
			subtotal += item.Price * float64(item.Quantity)
		}
	}

	// Use existing tax rate or default
	taxRate := order.TaxRate
	if taxRate == 0 {
		taxRate = 10
	}
	tax := subtotal * (taxRate / 100)

	serviceChargeRate := order.ServiceChargeRate
	serviceCharge := subtotal * (serviceChargeRate / 100)

	order.Subtotal = subtotal
	order.Tax = tax
	order.ServiceCharge = serviceCharge
	order.Total = subtotal + tax + serviceCharge
}