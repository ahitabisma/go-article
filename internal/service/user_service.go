package service

import (
	"errors"
	"go-article/internal/entity"
	"go-article/internal/repository"
)

type UserService interface {
	GetUserByID(userID uint) (entity.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

// GetUserByID implements UserService.
func (u *userService) GetUserByID(userID uint) (entity.User, error) {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return user, err
	}
	if user.ID == 0 {
		return user, errors.New("user not found")
	}
	return user, nil
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}
