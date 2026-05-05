package services

import (
	"errors"

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
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	if !user.Active {
		return nil, errors.New("account is deactivated")
	}

	// Update last login
	s.userRepo.UpdateLastLogin(user.ID)

	// Generate JWT
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

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Set default role if not specified
	role := req.Role
	if role == "" {
		role = models.RoleCashier
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         role,
		Active:       true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate JWT
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
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	if !utils.CheckPasswordHash(oldPassword, user.PasswordHash) {
		return errors.New("incorrect old password")
	}

	newHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.userRepo.Update(userID, map[string]interface{}{"passwordHash": newHash})
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
