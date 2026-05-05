package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name          string             `bson:"name" json:"name"`
	Icon          string             `bson:"icon" json:"icon"`
	BgColor       string             `bson:"bgColor" json:"bgColor"`
	SortOrder     int                `bson:"sortOrder" json:"sortOrder"`
	ShowInKitchen bool               `bson:"showInKitchen" json:"showInKitchen"`
	ShowInMenu    bool               `bson:"showInMenu" json:"showInMenu"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}