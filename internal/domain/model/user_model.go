package model

import (
	"time"
)

type User struct {
	ID         uint64     `gorm:"primaryKey;autoIncrement"`
	Name       string     `gorm:"type:varchar(255);not null"`
	Email      string     `gorm:"type:varchar(255);unique;not null;index:idx_users_email"`
	Password   string     `gorm:"type:varchar(255);not null"`
	Avatar     *string    `gorm:"type:varchar(512)"`
	VerifiedAt *time.Time `gorm:"column:verify_at"`
	Roles      []Role     `gorm:"many2many:user_role;"`
	CreatedAt  time.Time  `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt  time.Time  `gorm:"type:timestamp;default:current_timestamp on update current_timestamp"`
	DeletedAt  *time.Time `gorm:"index"`
}
