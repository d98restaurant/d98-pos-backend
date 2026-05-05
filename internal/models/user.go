package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleManager  UserRole = "manager"
	RoleCashier  UserRole = "cashier"
	RolePOS      UserRole = "pos"
	RoleKitchen  UserRole = "kitchen"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Username     string             `bson:"username" json:"username"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"passwordHash" json:"-"`
	Role         UserRole           `bson:"role" json:"role"`
	Active       bool               `bson:"active" json:"active"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
	LastLogin    *time.Time         `bson:"lastLogin,omitempty" json:"lastLogin,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string   `json:"username" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	Role     UserRole `json:"role"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	Role   UserRole `json:"role"`
	Active bool     `json:"active"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID        string    `json:"_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      UserRole  `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	LastLogin *time.Time `json:"lastLogin,omitempty"`
}