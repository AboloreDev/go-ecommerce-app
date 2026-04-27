package services

import (
	"github.com/aboloredev/armory/internal/dto"
	"github.com/aboloredev/armory/internal/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) GetUserProfileService(userID uint) (*dto.UserResponse, error) {
	var user models.User
	err := s.db.Where("id = ? AND is_active = ?", userID, true).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		Phone:     user.PhoneNumber,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *UserService) GetAllUsersService(page int, pageSize int) ([]*dto.UserResponse, int64, error) {
	var users []models.User
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	s.db.Model(&models.User{}).Count(&total)

	err := s.db.Where("role = ? AND is_active = ?", models.UserRoleCustomer, true).Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	response := make([]*dto.UserResponse, 0, len(users))

	for _, user := range users {
		response = append(response, &dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
			Phone:     user.PhoneNumber,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return response, total, nil
}

func (s *UserService) UpdateProfileService(userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	var user models.User
	err := s.db.Where("id = ? AND is_active = ?", userID, true).First(&user).Error
	if err != nil {
		return nil, err
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.PhoneNumber = req.Phone

	err = s.db.Save(&user).Error
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		Phone:     user.PhoneNumber,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
