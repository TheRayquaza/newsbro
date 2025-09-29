package services

import (
	"gorm.io/gorm"
	"repo_account/src/data/models"
)

type UserService struct {
	Db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{Db: db}
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := s.Db.First(&user, id).Error
	return &user, err
}

func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) (*models.User, error) {
	var user models.User
	if err := s.Db.First(&user, id).Error; err != nil {
		return nil, err
	}

	if role, ok := updates["role"].(string); ok {
		if role == "admin" && user.Role != "admin" {
			// only admin can assign admin role
			return nil, gorm.ErrInvalidData
		} else if role != "admin" && role != "user" {
			// only user or admin now
			return nil, gorm.ErrInvalidData
		}
	}

	if err := s.Db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetAllUsers(limit, offset int) ([]models.User, error) {
	var users []models.User
	err := s.Db.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}
