package repository

import (
	"go-article/internal/entity"
	"log"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user entity.User) (entity.User, error)
	FindByEmail(email string) (entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// Create implements UserRepository.
func (u *userRepository) Create(user entity.User) (entity.User, error) {
	err := u.db.Create(&user).Error
	if err != nil {
		log.Println("[UserRepository] Create:", err)
		return user, err
	}
	return user, nil
}

// FindByEmail implements UserRepository.
func (u *userRepository) FindByEmail(email string) (entity.User, error) {
	var user entity.User
	err := u.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Println("[UserRepository] FindByEmail:", err)
		return user, err
	}
	return user, nil

}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
