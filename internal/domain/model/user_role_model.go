package model

import "time"

type UserRole struct {
	UserID    uint64     `gorm:"primaryKey"`
	RoleID    uint64     `gorm:"primaryKey"`
	CreatedAt time.Time  `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt *time.Time `gorm:"type:timestamp;default:current_timestamp on update current_timestamp"`
}

func (UserRole) TableName() string {
	return "user_role"
}
