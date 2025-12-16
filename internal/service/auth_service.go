package service

import (
	"errors"
	"go-article/internal/domain/entity"
	"go-article/internal/handler/request"
	"go-article/internal/repository"
	"go-article/pkg/utils"
	"log"
	"strings"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(request request.RegisterRequest) (entity.UserEntity, error)
	Login(request request.LoginRequest) (entity.UserEntity, string, error)
	Profile(userID uint64) (entity.UserEntity, error)
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRepo,
	}
}

// Profile implements AuthService.
func (a *authService) Profile(userID uint64) (entity.UserEntity, error) {
	user, err := a.userRepository.FindByID(userID)
	if err != nil {
		log.Println("Error fetching user in Profile:", err)
		return user, err
	}
	return user, nil
}

// Login implements AuthService.
func (a *authService) Login(request request.LoginRequest) (entity.UserEntity, string, error) {
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
func (a *authService) Register(request request.RegisterRequest) (entity.UserEntity, error) {
	user := entity.UserEntity{}
	user.Name = request.Name
	user.Email = request.Email

	// Hash password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return user, err
	}

	user.Password = hashedPassword

	// Validasi role IDs ada di database
	if len(request.RoleIDs) > 0 {
		roles, err := a.userRepository.GetRolesByIDs(request.RoleIDs)
		if err != nil {
			log.Println("Error fetching roles:", err)
			return user, err
		}

		// Cek apakah jumlah role yang ditemukan sama dengan request
		if len(roles) != len(request.RoleIDs) {
			return user, errors.New("role IDs are invalid")
		}
	}

	newUser, err := a.userRepository.CreateWithRoles(user, request.RoleIDs)
	if err != nil {
		// Cek apakah error adalah duplicate key constraint
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return user, errors.New("email already registered")
		}
		return user, err
	}

	return newUser, nil
}
