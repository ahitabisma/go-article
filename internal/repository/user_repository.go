package repository

import (
	"go-article/internal/domain/entity"
	"go-article/internal/domain/model"
	"log"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user entity.UserEntity) (entity.UserEntity, error)
	CreateWithRoles(user entity.UserEntity, roleIDs []uint64) (entity.UserEntity, error)
	FindByEmail(email string) (entity.UserEntity, error)
	FindByID(id uint64) (entity.UserEntity, error)
	GetRolesByIDs(roleIDs []uint64) ([]model.Role, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// CreateWithRoles implements UserRepository.
func (u *userRepository) CreateWithRoles(user entity.UserEntity, roleIDs []uint64) (entity.UserEntity, error) {
	// Convert UserEntity to User model
	userModel := model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Avatar:   user.Avatar,
	}

	// Create user model first
	err := u.db.Create(&userModel).Error
	if err != nil {
		log.Println("[UserRepository] CreateWithRoles:", err)
		return user, err
	}

	// Associate roles if provided
	if len(roleIDs) > 0 {
		var roles []model.Role
		err = u.db.Where("id IN ?", roleIDs).Find(&roles).Error
		if err != nil {
			log.Println("[UserRepository] CreateWithRoles - fetching roles:", err)
			return user, err
		}

		err = u.db.Model(&userModel).Association("Roles").Append(roles)
		if err != nil {
			log.Println("[UserRepository] CreateWithRoles - associating roles:", err)
			return user, err
		}
	}

	// Convert back to UserEntity
	result := entity.UserEntity{
		ID:       userModel.ID,
		Name:     userModel.Name,
		Email:    userModel.Email,
		Password: userModel.Password,
		Avatar:   userModel.Avatar,
		Roles:    userModel.Roles,
	}

	return result, nil
}

// FindByID implements UserRepository.
func (u *userRepository) FindByID(id uint64) (entity.UserEntity, error) {
	var user model.User
	err := u.db.Where("id = ?", id).Preload("Roles").First(&user).Error
	if err != nil {
		log.Println("[UserRepository] FindByID:", err)
		return entity.UserEntity{}, err
	}

	return entity.UserEntity{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Roles: user.Roles,
	}, nil
}

// Create implements UserRepository.
func (u *userRepository) Create(user entity.UserEntity) (entity.UserEntity, error) {
	// Convert UserEntity to User model
	userModel := model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Avatar:   user.Avatar,
	}

	err := u.db.Create(&userModel).Error
	if err != nil {
		log.Println("[UserRepository] Create:", err)
		return user, err
	}

	// Convert back to UserEntity
	result := entity.UserEntity{
		ID:       userModel.ID,
		Name:     userModel.Name,
		Email:    userModel.Email,
		Password: userModel.Password,
		Avatar:   userModel.Avatar,
		Roles:    userModel.Roles,
	}

	return result, nil
}

// FindByEmail implements UserRepository.
func (u *userRepository) FindByEmail(email string) (entity.UserEntity, error) {
	var user model.User
	err := u.db.Where("email = ?", email).Preload("Roles").First(&user).Error
	if err != nil {
		log.Println("[UserRepository] FindByEmail:", err)
		return entity.UserEntity{}, err
	}

	return entity.UserEntity{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Avatar:   user.Avatar,
		Roles:    user.Roles,
	}, nil
}

// GetRolesByIDs implements UserRepository.
func (u *userRepository) GetRolesByIDs(roleIDs []uint64) ([]model.Role, error) {
	var roles []model.Role
	err := u.db.Where("id IN ?", roleIDs).Find(&roles).Error
	if err != nil {
		log.Println("[UserRepository] GetRolesByIDs:", err)
		return roles, err
	}
	return roles, nil
}
