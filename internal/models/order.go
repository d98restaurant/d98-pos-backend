package models

import (
	"time"
)

type OrderStatus string

const (
	OrderStatusPending        OrderStatus = "pending"
	OrderStatusAccepted       OrderStatus = "accepted"
	OrderStatusPreparing      OrderStatus = "preparing"
	OrderStatusReadyForBilling OrderStatus = "ready_for_billing"
	OrderStatusCompleted      OrderStatus = "completed"
	OrderStatusCancelled      OrderStatus = "cancelled"
	OrderStatusHold           OrderStatus = "hold"
)

type OrderType string

const (
	OrderTypeDineIn   OrderType = "dine-in"
	OrderTypeTakeaway OrderType = "takeaway"
	OrderTypeDelivery OrderType = "delivery"
)

type OrderItem struct {
	ID                  string     `json:"id"`
	Name                string     `json:"name"`
	Quantity            int        `json:"quantity"`
	Price               float64    `json:"price"`
	CategoryID          string     `json:"categoryId,omitempty"`
	CategoryName        string     `json:"categoryName,omitempty"`
	CategorySortOrder   int        `json:"categorySortOrder,omitempty"`
	SpecialInstructions string     `json:"specialInstructions,omitempty"`
	Status              string     `json:"status"`
	IsRemoved           bool       `json:"isRemoved,omitempty"`
	IsModified          bool       `json:"isModified,omitempty"`
	ModifiedAt          *time.Time `json:"modifiedAt,omitempty"`
	OldQuantity         int        `json:"oldQuantity,omitempty"`
	CancellationRequested bool     `json:"cancellationRequested,omitempty"`
	CancellationApproved  bool     `json:"cancellationApproved,omitempty"`
	CancellationReason    string   `json:"cancellationReason,omitempty"`
}

type CancellationRequest struct {
	OrderID         string    `json:"orderId"`
	OrderNumber     int       `json:"orderNumber"`
	ItemID          string    `json:"itemId"`
	ItemName        string    `json:"itemName"`
	Quantity        int       `json:"quantity"`
	Reason          string    `json:"reason"`
	RequestedAt     time.Time `json:"requestedAt"`
	OrderType       string    `json:"orderType"`
	TableNumber     int       `json:"tableNumber,omitempty"`
	DeliveryPlatform string   `json:"deliveryPlatform,omitempty"`
}

type PaymentDetails struct {
	Method        string         `json:"method"`
	Amount        float64        `json:"amount"`
	TransactionID string         `json:"transactionId,omitempty"`
	PaidAt        *time.Time     `json:"paidAt,omitempty"`
	Change        float64        `json:"change,omitempty"`
	Status        string         `json:"status"`
	DueDate       *time.Time     `json:"dueDate,omitempty"`
	CustomerName  string         `json:"customerName,omitempty"`
	CustomerPhone string         `json:"customerPhone,omitempty"`
	CustomerEmail string         `json:"customerEmail,omitempty"`
	CustomerID    string         `json:"customerId,omitempty"`
	Notes         string         `json:"notes,omitempty"`
	SplitDetails  []SplitPayment `json:"splitDetails,omitempty"`
}

type SplitPayment struct {
	Method string  `json:"method"`
	Amount float64 `json:"amount"`
}

type Customer struct {
	Name  string `json:"name"`
	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`
	ID    string `json:"id,omitempty"`
}

type Order struct {
	ID                  string          `json:"_id,omitempty"`
	OrderNumber         int             `json:"orderNumber"`
	DisplayOrderNumber  string          `json:"displayOrderNumber,omitempty"`
	BaseOrderNumber     int             `json:"baseOrderNumber,omitempty"`
	RunningNumber       int             `json:"runningNumber,omitempty"`
	Items               []OrderItem     `json:"items"`
	Subtotal            float64         `json:"subtotal"`
	Tax                 float64         `json:"tax"`
	TaxRate             float64         `json:"taxRate"`
	ServiceCharge       float64         `json:"serviceCharge"`
	ServiceChargeRate   float64         `json:"serviceChargeRate"`
	Total               float64         `json:"total"`
	OrderType           OrderType       `json:"orderType"`
	TableNumber         int             `json:"tableNumber,omitempty"`
	DeliveryPlatform    string          `json:"deliveryPlatform,omitempty"`
	DeliveryAddress     string          `json:"deliveryAddress,omitempty"`
	Status              OrderStatus     `json:"status"`
	Customer            Customer        `json:"customer"`
	Payment             *PaymentDetails `json:"payment,omitempty"`
	CreatedAt           time.Time       `json:"createdAt"`
	UpdatedAt           time.Time       `json:"updatedAt"`
	CompletedAt         *time.Time      `json:"completedAt,omitempty"`
	IsAdditionalOrder   bool            `json:"isAdditionalOrder,omitempty"`
	ParentOrderID       string          `json:"parentOrderId,omitempty"`
	TableSessionID      string          `json:"tableSessionId,omitempty"`
	HasModifications    bool            `json:"hasModifications,omitempty"`
}
