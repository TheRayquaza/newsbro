package services

import (
	"gorm.io/gorm"
	"repo_account/src/data/models"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, id).Error
	return &user, err
}

func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetAllUsers(limit, offset int) ([]models.User, error) {
	var users []models.User
	err := s.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}
