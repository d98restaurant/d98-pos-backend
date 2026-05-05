package services

import (
	"errors"
	"fmt"
	"log"

	"pos-backend/config"
	"pos-backend/internal/models"
	"pos-backend/internal/repository"
	"pos-backend/internal/utils"

	"golang.org/x/crypto/bcrypt"
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
	log.Printf("🔐 Login attempt: username=%s", username)
	
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
	
	// Direct bcrypt verification
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Printf("❌ Password mismatch for user: %s - Error: %v", username, err)
		return nil, errors.New("invalid username or password")
	}

	log.Printf("✅ Password verified successfully")

	if !user.Active {
		log.Printf("❌ Account deactivated: %s", username)
		return nil, errors.New("account is deactivated")
	}

	// Update last login
	go s.userRepo.UpdateLastLogin(user.ID)

	// Generate JWT
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

	log.Printf("✅ Login successful for user: %s", username)
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
	log.Printf("📝 Registering new user: %s", req.Username)
	
	// Check if username exists
	existingUser, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email exists
	existingEmail, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingEmail != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password using bcrypt with cost 10
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("❌ Password hashing error: %v", err)
		return nil, errors.New("failed to hash password")
	}
	passwordHash := string(hashedBytes)
	
	log.Printf("✅ Password hashed successfully - Hash length: %d", len(passwordHash))
	log.Printf("   Hash prefix: %s", passwordHash[:3])

	// Set default role if not specified
	role := req.Role
	if role == "" {
		role = models.RoleCashier
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         role,
		Active:       true,
	}

	if err := s.userRepo.Create(user); err != nil {
		log.Printf("❌ User creation error: %v", err)
		return nil, err
	}

	log.Printf("✅ User created successfully: ID=%s, Username=%s", user.ID, user.Username)

	// Generate JWT
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
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("incorrect old password")
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.Update(userID, map[string]interface{}{"passwordHash": string(newHash)})
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
