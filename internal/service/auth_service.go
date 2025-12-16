package service

import (
	"errors"
	"go-article/internal/entity"
	"go-article/internal/handler/request"
	"go-article/internal/repository"
	"go-article/pkg/utils"
	"log"
	"strings"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(request request.RegisterRequest) (entity.User, error)
	Login(request request.LoginRequest) (entity.User, string, error)
}

type authService struct {
	userRepository repository.UserRepository
}

// Login implements AuthService.
func (a *authService) Login(request request.LoginRequest) (entity.User, string, error) {
	// Cari user berdasarkan email
	user, err := a.userRepository.FindByEmail(request.Email)
	if err != nil {
		return user, "", errors.New("invalid email or password")
	}

	// Periksa apakah user ditemukan
	if user.ID == 0 {
		return user, "", errors.New("invalid email or password")
	}

	// Verifikasi password
	if !utils.CheckPasswordHash(request.Password, user.Password) {
		return user, "", errors.New("invalid email or password")
	}

	// Generate token JWT
	token, err := utils.GenerateToken(uint(user.ID))
	if err != nil {
		return user, "", err
	}

	return user, token, nil
}

// Register implements AuthService.
func (a *authService) Register(request request.RegisterRequest) (entity.User, error) {
	user := entity.User{}
	user.Name = request.Name
	user.Email = request.Email

	// Hash password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return user, err
	}

	user.Password = hashedPassword

	newUser, err := a.userRepository.Create(user)
	if err != nil {
		// Cek apakah error adalah duplicate key constraint
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return user, errors.New("email already registered")
		}
		return user, err
	}

	return newUser, nil
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRepo,
	}
}
