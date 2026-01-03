package model

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/datatypes"
)

type Portfolio struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Title       string         `gorm:"type:varchar(255)" json:"title"`
	Bio         string         `gorm:"type:text" json:"bio"`
	Email       string         `gorm:"type:varchar(255);index" json:"email"`
	SocialLinks datatypes.JSON `gorm:"type:jsonb" json:"social_links"`
	Settings    datatypes.JSON `gorm:"type:jsonb" json:"settings"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (p *Portfolio) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (p *Portfolio) TableName() string {
	return "portfolio"
}

