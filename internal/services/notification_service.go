package services

import (
	"context"
	"encoding/base64"
	"log"

	"pos-backend/internal/models"
)

var vapidPublicKey = "BLiNlzBDszn00bfL0w25zye0nz725AzWDHT-9RM-pksxzBbMWiO9W8cVTZdL1W5ouLIAcTMKyFkN4eXq4vAb_Kg"

type NotificationService struct {
	subscriptions map[string]map[string]interface{}
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		subscriptions: make(map[string]map[string]interface{}),
	}
}

func GetVAPIDPublicKey() string {
	return vapidPublicKey
}

func (s *NotificationService) Subscribe(ctx context.Context, userID, role string, subscription map[string]interface{}) error {
	s.subscriptions[userID] = subscription
	log.Printf("User %s (%s) subscribed to notifications", userID, role)
	return nil
}

func (s *NotificationService) Unsubscribe(ctx context.Context, userID string) error {
	delete(s.subscriptions, userID)
	return nil
}

func (s *NotificationService) GetKitchenSubscriberCount() int {
	return len(s.subscriptions)
}

func (s *NotificationService) SendOrderNotification(order *models.Order) {
	// In a real implementation, this would send push notifications
	// For now, just log
	log.Printf("📢 Notification: New Order #%d (%s)", order.OrderNumber, order.OrderType)
}

func (s *NotificationService) SendNotificationToUser(userID string, title, body string) {
	log.Printf("🔔 Notification to %s: %s - %s", userID, title, body)
}