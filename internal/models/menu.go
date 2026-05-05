package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MenuItem struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	FullName    string             `bson:"fullName,omitempty" json:"fullName,omitempty"`
	Price       float64            `bson:"price" json:"price"`
	Category    primitive.ObjectID `bson:"category" json:"category"`
	CategoryName string             `bson:"categoryName,omitempty" json:"categoryName,omitempty"`
	PrepTime    int                `bson:"prepTime" json:"prepTime"`
	Available   bool               `bson:"available" json:"available"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	ImageURL    string             `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}