package models

import (
	"time"
)

type Category struct {
	ID            string    `json:"_id,omitempty"`
	Name          string    `json:"name"`
	Icon          string    `json:"icon"`
	BgColor       string    `json:"bgColor"`
	SortOrder     int       `json:"sortOrder"`
	ShowInKitchen bool      `json:"showInKitchen"`
	ShowInMenu    bool      `json:"showInMenu"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
