package services

import (
	"context"
	"log"

	"pos-backend/config"
	"pos-backend/internal/models"
)

type NotificationService struct {
	config        *config.Config
	subscriptions map[string]map[string]interface{}
}

func NewNotificationService(cfg *config.Config) *NotificationService {
	return &NotificationService{
		config:        cfg,
		subscriptions: make(map[string]map[string]interface{}),
	}
}

func (s *NotificationService) GetVAPIDPublicKey() string {
	return s.config.VAPIDPublicKey
}

func (s *NotificationService) Subscribe(ctx context.Context, userID, role string, subscription map[string]interface{}) error {
	s.subscriptions[userID] = subscription
	log.Printf("✅ User %s (%s) subscribed to push notifications", userID, role)
	return nil
}

func (s *NotificationService) Unsubscribe(ctx context.Context, userID string) error {
	delete(s.subscriptions, userID)
	log.Printf("❌ User %s unsubscribed from push notifications", userID)
	return nil
}

func (s *NotificationService) GetKitchenSubscriberCount() int {
	return len(s.subscriptions)
}

func (s *NotificationService) SendOrderNotification(order *models.Order) {
	for userID, sub := range s.subscriptions {
		log.Printf("🔔 Sending notification to %s: New Order #%d (%s)", 
			userID, order.OrderNumber, order.OrderType)
	}
}

func (s *NotificationService) SendNotificationToUser(userID, title, body string) {
	if _, exists := s.subscriptions[userID]; exists {
		log.Printf("📢 Sending to %s: %s - %s", userID, title, body)
	}
}
