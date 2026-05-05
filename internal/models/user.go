package models

import (
	"time"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleManager UserRole = "manager"
	RoleCashier UserRole = "cashier"
	RolePOS     UserRole = "pos"
	RoleKitchen UserRole = "kitchen"
)

type User struct {
	ID           string     `json:"_id" bson:"_id,omitempty"`
	Username     string     `json:"username" bson:"username"`
	Email        string     `json:"email" bson:"email"`
	PasswordHash string     `json:"-" bson:"passwordHash"` // Never send password in JSON
	Role         UserRole   `json:"role" bson:"role"`
	Active       bool       `json:"active" bson:"active"`
	CreatedAt    time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt" bson:"updatedAt"`
	LastLogin    *time.Time `json:"lastLogin,omitempty" bson:"lastLogin,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string   `json:"username" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=3"`
	Role     UserRole `json:"role"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=3"`
}

type UpdateUserRequest struct {
	Role   UserRole `json:"role"`
	Active bool     `json:"active"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID        string     `json:"_id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Role      UserRole   `json:"role"`
	Active    bool       `json:"active"`
	CreatedAt time.Time  `json:"createdAt"`
	LastLogin *time.Time `json:"lastLogin,omitempty"`
}
