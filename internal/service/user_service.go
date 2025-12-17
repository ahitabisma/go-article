package service

import (
	"errors"
	"go-article/internal/domain/entity"
	"go-article/internal/repository"
)

type UserService interface {
	GetUserByID(userID uint64) (*entity.UserEntity, error)
}

type userService struct {
	userRepo repository.UserRepository
}

// GetUserByID implements UserService.
func (u *userService) GetUserByID(userID uint64) (*entity.UserEntity, error) {
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
