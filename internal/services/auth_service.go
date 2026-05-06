package services

import (
	"errors"
	"log"

	"pos-backend/config"
	"pos-backend/internal/models"
	"pos-backend/internal/repository"
	"pos-backend/internal/utils"
)

type AuthService struct {
	userRepo *repository.UserRepository
	config   *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
	}
}

func (s *AuthService) Login(username, password string) (*models.AuthResponse, error) {
	log.Printf("🔐 Login attempt: %s", username)
	
	// Find user by username
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		log.Printf("❌ Database error: %v", err)
		return nil, errors.New("invalid username or password")
	}
	
	if user == nil {
		log.Printf("❌ User not found: %s", username)
		return nil, errors.New("invalid username or password")
	}
	
	log.Printf("✅ User found: %s", user.Username)
	
	// Check password (plain text comparison)
	if user.PasswordHash != password {
		log.Printf("❌ Password mismatch for: %s", username)
		return nil, errors.New("invalid username or password")
	}
	
	log.Printf("✅ Password verified for: %s", username)
	
	// Check if account is active
	if !user.Active {
		log.Printf("❌ Account deactivated: %s", username)
		return nil, errors.New("account is deactivated. Please contact admin.")
	}
	
	// Update last login
	go s.userRepo.UpdateLastLogin(user.ID)
	
	// Generate JWT token
	token, err := utils.GenerateJWT(
		user.ID,
		user.Username,
		string(user.Role),
		s.config.JWTSecret,
		s.config.JWTExpiry,
	)
	if err != nil {
		log.Printf("❌ Token generation error: %v", err)
		return nil, err
	}
	
	log.Printf("✅ Login successful for: %s (Role: %s)", username, user.Role)
	
	return &models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			Active:    user.Active,
			CreatedAt: user.CreatedAt,
			LastLogin: user.LastLogin,
		},
	}, nil
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	log.Printf("📝 Register attempt: %s", req.Username)
	
	// Check if username exists
	existing, _ := s.userRepo.FindByUsername(req.Username)
	if existing != nil {
		log.Printf("❌ Username already exists: %s", req.Username)
		return nil, errors.New("username already exists")
	}
	
	// Check if email exists
	existingEmail, _ := s.userRepo.FindByEmail(req.Email)
	if existingEmail != nil {
		log.Printf("❌ Email already registered: %s", req.Email)
		return nil, errors.New("email already registered")
	}
	
	// Set role (default to cashier for security)
	role := req.Role
	if role == "" {
		role = models.RoleCashier
	}
	
	// Create user with plain text password
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password,
		Role:         role,
		Active:       true,
	}
	
	if err := s.userRepo.Create(user); err != nil {
		log.Printf("❌ User creation error: %v", err)
		return nil, err
	}
	
	log.Printf("✅ User created: %s (ID: %s, Role: %s)", user.Username, user.ID, user.Role)
	
	// Generate token
	token, err := utils.GenerateJWT(
		user.ID,
		user.Username,
		string(user.Role),
		s.config.JWTSecret,
		s.config.JWTExpiry,
	)
	if err != nil {
		return nil, err
	}
	
	return &models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			Active:    user.Active,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *AuthService) ChangePassword(userID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	
	if user.PasswordHash != oldPassword {
		return errors.New("incorrect old password")
	}
	
	return s.userRepo.Update(userID, map[string]interface{}{"passwordHash": newPassword})
}

func (s *AuthService) GetUsers() ([]models.UserResponse, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}
	
	responses := make([]models.UserResponse, 0, len(users))
	for _, user := range users {
		if user.ID == "" || user.Username == "" {
			continue
		}
		responses = append(responses, models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			Active:    user.Active,
			CreatedAt: user.CreatedAt,
			LastLogin: user.LastLogin,
		})
	}
	
	log.Printf("📋 Returning %d users", len(responses))
	return responses, nil
}

func (s *AuthService) UpdateUser(userID string, req *models.UpdateUserRequest) error {
	updates := map[string]interface{}{
		"role":   req.Role,
		"active": req.Active,
	}
	return s.userRepo.Update(userID, updates)
}

func (s *AuthService) DeleteUser(userID string) error {
	return s.userRepo.Delete(userID)
}
