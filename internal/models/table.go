package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TableStatus string

const (
	TableStatusAvailable TableStatus = "available"
	TableStatusOccupied  TableStatus = "occupied"
	TableStatusReserved  TableStatus = "reserved"
)

type Table struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	TableNumber      int                `bson:"tableNumber" json:"tableNumber"`
	Capacity         int                `bson:"capacity" json:"capacity"`
	Status           TableStatus        `bson:"status" json:"status"`
	CurrentSessionID string             `bson:"currentSessionId,omitempty" json:"currentSessionId,omitempty"`
	RunningOrderCount int               `bson:"runningOrderCount" json:"runningOrderCount"`
	TotalRunningAmount float64          `bson:"totalRunningAmount" json:"totalRunningAmount"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}