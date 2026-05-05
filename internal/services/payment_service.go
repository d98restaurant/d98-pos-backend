package services

import "pos-backend/config"

type PaymentService struct {
	config *config.Config
}

func NewPaymentService(cfg *config.Config) *PaymentService {
	return &PaymentService{config: cfg}
}

func (s *PaymentService) GetKeyID() string {
	return s.config.RazorpayKeyID
}

func (s *PaymentService) GetKeySecret() string {
	return s.config.RazorpayKeySecret
}