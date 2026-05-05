package models

import (
	"time"
)

type TableStatus string

const (
	TableStatusAvailable TableStatus = "available"
	TableStatusOccupied  TableStatus = "occupied"
	TableStatusReserved  TableStatus = "reserved"
)

type Table struct {
	ID                string       `json:"_id,omitempty"`
	TableNumber       int          `json:"tableNumber"`
	Capacity          int          `json:"capacity"`
	Status            TableStatus  `json:"status"`
	CurrentSessionID  string       `json:"currentSessionId,omitempty"`
	RunningOrderCount int          `json:"runningOrderCount"`
	TotalRunningAmount float64     `json:"totalRunningAmount"`
	CreatedAt         time.Time    `json:"createdAt"`
	UpdatedAt         time.Time    `json:"updatedAt"`
}
