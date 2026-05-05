package services

import (
	"errors"
	"fmt"
	"time"

	"pos-backend/internal/models"
	"pos-backend/internal/repository"
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

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.FindAll()
}

func (s *OrderService) GetOrderByID(id string) (*models.Order, error) {
	return s.orderRepo.FindByID(id)
}

func (s *OrderService) CreateOrder(order *models.Order) (*models.Order, error) {
	orderNumber, err := s.orderRepo.GetNextOrderNumber()
	if err != nil {
		orderNumber = int(time.Now().UnixNano() % 1000000)
	}
	order.OrderNumber = orderNumber

	if order.IsAdditionalOrder && order.TableNumber > 0 {
		table, err := s.tableRepo.FindByNumber(order.TableNumber)
		if err == nil && table != nil && table.RunningOrderCount > 0 {
			order.RunningNumber = table.RunningOrderCount + 1
			order.BaseOrderNumber = order.OrderNumber
			order.DisplayOrderNumber = fmt.Sprintf("%d-%d", order.OrderNumber, order.RunningNumber)
		}
	}

	order.Status = models.OrderStatusPending

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	if order.OrderType == models.OrderTypeDineIn && order.TableNumber > 0 {
		s.tableRepo.IncrementRunningCount(order.TableNumber, order.Total)
	}

	s.notification.SendOrderNotification(order)

	return order, nil
}

func (s *OrderService) UpdateOrderStatus(id string, status models.OrderStatus) (*models.Order, error) {
	if err := s.orderRepo.UpdateStatus(id, status); err != nil {
		return nil, err
	}
	return s.orderRepo.FindByID(id)
}

func (s *OrderService) CompletePayment(id, paymentMethod string, paymentDetails map[string]interface{}) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(id)
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

	if err := s.orderRepo.UpdateOrderPayment(id, paymentMethod, paymentDetails); err != nil {
		return nil, err
	}

	if order.OrderType == models.OrderTypeDineIn && order.TableNumber > 0 {
		s.tableRepo.DecrementRunningCount(order.TableNumber, order.Total)
	}

	return s.orderRepo.FindByID(id)
}

func (s *OrderService) AddItemToOrder(orderID string, item *models.OrderItem) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	// Rest of the implementation...
	return order, nil
}

func (s *OrderService) UpdateItemQuantity(orderID, itemID string, quantity int) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	// Rest of the implementation...
	return order, nil
}

func (s *OrderService) RemoveItemFromOrder(orderID, itemID string) (*models.Order, error) {
	return s.UpdateItemQuantity(orderID, itemID, 0)
}
