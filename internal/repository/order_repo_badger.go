package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"pos-backend/internal/models"

	"github.com/dgraph-io/badger/v4"
)

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (r *OrderRepository) Create(order *models.Order) error {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	
	orderNumber, _ := GetNextSequence("order_number")
	order.OrderNumber = int(orderNumber)
	order.ID = fmt.Sprintf("%d", orderNumber)
	
	return SaveJSON("order:"+order.ID, order)
}

func (r *OrderRepository) FindByID(id string) (*models.Order, error) {
	var order models.Order
	err := GetJSON("order:"+id, &order)
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	return &order, err
}

func (r *OrderRepository) FindAll() ([]models.Order, error) {
	var orders []models.Order
	err := DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := []byte("order:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var order models.Order
				if err := json.Unmarshal(val, &order); err != nil {
					return err
				}
				orders = append(orders, order)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return orders, err
}

func (r *OrderRepository) UpdateStatus(id string, status models.OrderStatus) error {
	order, err := r.FindByID(id)
	if err != nil {
		return err
	}
	if order == nil {
		return fmt.Errorf("order not found")
	}
	order.Status = status
	order.UpdatedAt = time.Now()
	return SaveJSON("order:"+id, order)
}

func (r *OrderRepository) Update(id string, updates map[string]interface{}) error {
	order, err := r.FindByID(id)
	if err != nil {
		return err
	}
	if order == nil {
		return fmt.Errorf("order not found")
	}
	
	if status, ok := updates["status"]; ok {
		order.Status = status.(models.OrderStatus)
	}
	if payment, ok := updates["payment"]; ok {
		order.Payment = payment.(*models.PaymentDetails)
	}
	if items, ok := updates["items"]; ok {
		order.Items = items.([]models.OrderItem)
	}
	order.UpdatedAt = time.Now()
	
	return SaveJSON("order:"+id, order)
}

func (r *OrderRepository) UpdateOrderPayment(id string, paymentMethod string, paymentDetails map[string]interface{}) error {
	order, err := r.FindByID(id)
	if err != nil {
		return err
	}
	if order == nil {
		return fmt.Errorf("order not found")
	}
	
	order.Payment = &models.PaymentDetails{
		Method:        paymentMethod,
		Amount:        paymentDetails["amount"].(float64),
		TransactionID: paymentDetails["transactionId"].(string),
		Status:        "completed",
	}
	paidAt := time.Now()
	order.Payment.PaidAt = &paidAt
	order.Status = models.OrderStatusCompleted
	order.UpdatedAt = time.Now()
	completedAt := time.Now()
	order.CompletedAt = &completedAt
	
	return SaveJSON("order:"+id, order)
}

func (r *OrderRepository) GetNextOrderNumber() (int, error) {
	seq, err := GetNextSequence("order_number")
	return int(seq), err
}

func (r *OrderRepository) FindCreditCustomers() ([]map[string]interface{}, error) {
	// For BadgerDB, we'll implement this later if needed
	return []map[string]interface{}{}, nil
}
