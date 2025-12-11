package repository

import (
	"natsApi/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *models.UserMessage) error {
	result := r.DB.Create(user)
	return result.Error
}

func (r *UserRepository) UpdateUserRoleByUsername(username string, newRole string) error {
	result := r.DB.Model(&models.UserMessage{}).
		Where("username = ?", username).
		Update("role", newRole)

	return result.Error
}
