package entity

import "go-article/internal/domain/model"

type UserEntity struct {
	ID         uint64
	Name       string
	Email      string
	Password   string `json:"-"`
	Avatar     *string
	VerifiedAt *string
	Roles      []model.Role
}
