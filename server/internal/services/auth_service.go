package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/aboloredev/armory/internal/config"
	"github.com/aboloredev/armory/internal/dto"
	"github.com/aboloredev/armory/internal/models"
	"github.com/aboloredev/armory/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db     *gorm.DB
	config *config.Config
}

var ErrNotFound = errors.New("User Not found")
var ErrConflict = errors.New("This email is already in use")
var ErrTokeNotFoundOrExpired = errors.New("Token Not found or expired")
var ErrInvalidCredentials = errors.New("Invalid Credentials")
var ErrInvalidRefreshToken = errors.New("Invalid Refresh Token")

func NewAuthService(db *gorm.DB, config *config.Config) *AuthService {
	return &AuthService{
		db:     db,
		config: config,
	}
}

func (s *AuthService) UserRegistrationService(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	var existingUser models.User

	err := s.db.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		return nil, ErrConflict
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:       req.Email,
		Password:    hashedPassword,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.Phone,
		Role:        models.UserRoleCustomer,
	}

	err = s.db.Create(&user).Error
	if err != nil {
		return nil, err
	}

	cart := models.Cart{
		UserID: user.ID,
	}

	err = s.db.Create(&cart).Error
	if err != nil {
		fmt.Println("Failed to create cart")
	}

	return s.GenerateAuthResponse(&user)
}

func (s *AuthService) UserLoginService(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	var user models.User

	err := s.db.Where("email = ? AND is_active = ?", req.Email, true).First(&user).Error
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, ErrInvalidCredentials
	}

	return s.GenerateAuthResponse(&user)

}

func (s *AuthService) RefreshTokenService(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken, s.config.JWT.JWTSecret)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	var refreshToken models.RefreshToken
	err = s.db.Where("token = ? AND expires_at > ?", req.RefreshToken, time.Now()).First(&refreshToken).Error
	if err != nil {
		return nil, ErrTokeNotFoundOrExpired
	}

	var user models.User
	err = s.db.First(&user, claims.UserID).Error
	if err != nil {
		return nil, ErrNotFound
	}

	s.db.Delete(&refreshToken)

	return s.GenerateAuthResponse(&user)
}

func (s *AuthService) LogoutService(req *dto.RefreshTokenRequest) error {
	var refreshToken models.RefreshToken

	err := s.db.Where("token = ?", req.RefreshToken).Delete(&refreshToken).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GenerateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokenPair(&s.config.JWT, user.ID, user.Email, string(user.Role))

	if err != nil {
		return nil, err
	}

	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.config.JWT.JWTRefreshTokenExpiration),
	}

	s.db.Create(&refreshTokenModel)

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.PhoneNumber,
			IsActive:  user.IsActive,
			Role:      string(user.Role),
		},
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil

}
