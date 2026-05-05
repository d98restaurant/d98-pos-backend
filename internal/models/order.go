package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	OrderStatusPending       OrderStatus = "pending"
	OrderStatusAccepted      OrderStatus = "accepted"
	OrderStatusPreparing     OrderStatus = "preparing"
	OrderStatusReadyForBilling OrderStatus = "ready_for_billing"
	OrderStatusCompleted     OrderStatus = "completed"
	OrderStatusCancelled     OrderStatus = "cancelled"
	OrderStatusHold          OrderStatus = "hold"
)

type OrderType string

const (
	OrderTypeDineIn    OrderType = "dine-in"
	OrderTypeTakeaway  OrderType = "takeaway"
	OrderTypeDelivery  OrderType = "delivery"
)

type OrderItem struct {
	ID                  string    `bson:"id" json:"id"`
	Name                string    `bson:"name" json:"name"`
	Quantity            int       `bson:"quantity" json:"quantity"`
	Price               float64   `bson:"price" json:"price"`
	CategoryID          string    `bson:"categoryId,omitempty" json:"categoryId,omitempty"`
	CategoryName        string    `bson:"categoryName,omitempty" json:"categoryName,omitempty"`
	CategorySortOrder   int       `bson:"categorySortOrder,omitempty" json:"categorySortOrder,omitempty"`
	SpecialInstructions string    `bson:"specialInstructions,omitempty" json:"specialInstructions,omitempty"`
	Status              string    `bson:"status" json:"status"`
	IsRemoved           bool      `bson:"isRemoved,omitempty" json:"isRemoved,omitempty"`
	IsModified          bool      `bson:"isModified,omitempty" json:"isModified,omitempty"`
	ModifiedAt          *time.Time `bson:"modifiedAt,omitempty" json:"modifiedAt,omitempty"`
	OldQuantity         int       `bson:"oldQuantity,omitempty" json:"oldQuantity,omitempty"`
	CancellationRequested bool    `bson:"cancellationRequested,omitempty" json:"cancellationRequested,omitempty"`
	CancellationApproved  bool    `bson:"cancellationApproved,omitempty" json:"cancellationApproved,omitempty"`
	CancellationReason    string  `bson:"cancellationReason,omitempty" json:"cancellationReason,omitempty"`
}

type CancellationRequest struct {
	OrderID     primitive.ObjectID `bson:"orderId" json:"orderId"`
	OrderNumber int                `bson:"orderNumber" json:"orderNumber"`
	ItemID      string             `bson:"itemId" json:"itemId"`
	ItemName    string             `bson:"itemName" json:"itemName"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	Reason      string             `bson:"reason" json:"reason"`
	RequestedAt time.Time          `bson:"requestedAt" json:"requestedAt"`
	OrderType   string             `bson:"orderType" json:"orderType"`
	TableNumber int                `bson:"tableNumber,omitempty" json:"tableNumber,omitempty"`
	DeliveryPlatform string        `bson:"deliveryPlatform,omitempty" json:"deliveryPlatform,omitempty"`
}

type PaymentDetails struct {
	Method          string     `bson:"method" json:"method"`
	Amount          float64    `bson:"amount" json:"amount"`
	TransactionID   string     `bson:"transactionId,omitempty" json:"transactionId,omitempty"`
	PaidAt          *time.Time `bson:"paidAt,omitempty" json:"paidAt,omitempty"`
	Change          float64    `bson:"change,omitempty" json:"change,omitempty"`
	Status          string     `bson:"status" json:"status"`
	DueDate         *time.Time `bson:"dueDate,omitempty" json:"dueDate,omitempty"`
	CustomerName    string     `bson:"customerName,omitempty" json:"customerName,omitempty"`
	CustomerPhone   string     `bson:"customerPhone,omitempty" json:"customerPhone,omitempty"`
	CustomerEmail   string     `bson:"customerEmail,omitempty" json:"customerEmail,omitempty"`
	CustomerID      string     `bson:"customerId,omitempty" json:"customerId,omitempty"`
	Notes           string     `bson:"notes,omitempty" json:"notes,omitempty"`
	SplitDetails    []SplitPayment `bson:"splitDetails,omitempty" json:"splitDetails,omitempty"`
}

type SplitPayment struct {
	Method string  `bson:"method" json:"method"`
	Amount float64 `bson:"amount" json:"amount"`
}

type Customer struct {
	Name  string `bson:"name" json:"name"`
	Phone string `bson:"phone,omitempty" json:"phone,omitempty"`
	Email string `bson:"email,omitempty" json:"email,omitempty"`
	ID    string `bson:"id,omitempty" json:"id,omitempty"`
}

type Order struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	OrderNumber         int                `bson:"orderNumber" json:"orderNumber"`
	DisplayOrderNumber  string             `bson:"displayOrderNumber,omitempty" json:"displayOrderNumber,omitempty"`
	BaseOrderNumber     int                `bson:"baseOrderNumber,omitempty" json:"baseOrderNumber,omitempty"`
	RunningNumber       int                `bson:"runningNumber,omitempty" json:"runningNumber,omitempty"`
	Items               []OrderItem        `bson:"items" json:"items"`
	Subtotal            float64            `bson:"subtotal" json:"subtotal"`
	Tax                 float64            `bson:"tax" json:"tax"`
	TaxRate             float64            `bson:"taxRate" json:"taxRate"`
	ServiceCharge       float64            `bson:"serviceCharge" json:"serviceCharge"`
	ServiceChargeRate   float64            `bson:"serviceChargeRate" json:"serviceChargeRate"`
	Total               float64            `bson:"total" json:"total"`
	OrderType           OrderType          `bson:"orderType" json:"orderType"`
	TableNumber         int                `bson:"tableNumber,omitempty" json:"tableNumber,omitempty"`
	DeliveryPlatform    string             `bson:"deliveryPlatform,omitempty" json:"deliveryPlatform,omitempty"`
	DeliveryAddress     string             `bson:"deliveryAddress,omitempty" json:"deliveryAddress,omitempty"`
	Status              OrderStatus        `bson:"status" json:"status"`
	Customer            Customer           `bson:"customer" json:"customer"`
	Payment             *PaymentDetails    `bson:"payment,omitempty" json:"payment,omitempty"`
	CreatedAt           time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt           time.Time          `bson:"updatedAt" json:"updatedAt"`
	CompletedAt         *time.Time         `bson:"completedAt,omitempty" json:"completedAt,omitempty"`
	IsAdditionalOrder   bool               `bson:"isAdditionalOrder,omitempty" json:"isAdditionalOrder,omitempty"`
	ParentOrderID       primitive.ObjectID `bson:"parentOrderId,omitempty" json:"parentOrderId,omitempty"`
	TableSessionID      string             `bson:"tableSessionId,omitempty" json:"tableSessionId,omitempty"`
	HasModifications    bool               `bson:"hasModifications,omitempty" json:"hasModifications,omitempty"`
}