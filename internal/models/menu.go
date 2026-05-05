package models

import (
	"time"
)

type MenuItem struct {
	ID           string    `json:"_id,omitempty"`
	Name         string    `json:"name"`
	FullName     string    `json:"fullName,omitempty"`
	Price        float64   `json:"price"`
	Category     string    `json:"category"`
	CategoryName string    `json:"categoryName,omitempty"`
	PrepTime     int       `json:"prepTime"`
	Available    bool      `json:"available"`
	Description  string    `json:"description,omitempty"`
	ImageURL     string    `json:"imageUrl,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
