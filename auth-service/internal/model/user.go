package model

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email         string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash  string    `gorm:"not null" json:"-"`
	Name          string    `json:"name"`
	Role          string    `gorm:"default:user" json:"role"`
	EmailVerified bool      `gorm:"default:false" json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (u *User) TableName() string {
	return "users"
}

