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
	log.Printf("🔐 Login: %s", username)
	
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return nil, errors.New("invalid username or password")
	}
	if user == nil {
		log.Printf("User not found: %s", username)
		return nil, errors.New("invalid username or password")
	}
	
	log.Printf("Found user: %s, stored password: '%s'", user.Username, user.PasswordHash)
	log.Printf("Provided password: '%s'", password)
	
	// Direct comparison
	if user.PasswordHash != password {
		log.Printf("Password mismatch!")
		return nil, errors.New("invalid username or password")
	}
	
	log.Printf("✅ Login successful: %s", username)
	
	// Update last login
	go s.userRepo.UpdateLastLogin(user.ID)
	
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
			LastLogin: user.LastLogin,
		},
	}, nil
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	log.Printf("📝 Register: %s with password: '%s'", req.Username, req.Password)
	
	// Check if username exists
	existing, _ := s.userRepo.FindByUsername(req.Username)
	if existing != nil {
		return nil, errors.New("username already exists")
	}
	
	// Check if email exists
	existingEmail, _ := s.userRepo.FindByEmail(req.Email)
	if existingEmail != nil {
		return nil, errors.New("email already registered")
	}
	
	// Set role
	role := req.Role
	if role == "" {
		role = models.RoleCashier
	}
	
	// Create user with plain text password
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password, // Store exactly as provided
		Role:         role,
		Active:       true,
	}
	
	if err := s.userRepo.Create(user); err != nil {
		log.Printf("Create error: %v", err)
		return nil, err
	}
	
	log.Printf("✅ User created: %s (ID: %s) with password: '%s'", user.Username, user.ID, user.PasswordHash)
	
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
	
	responses := make([]models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			Active:    user.Active,
			CreatedAt: user.CreatedAt,
			LastLogin: user.LastLogin,
		}
	}
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
