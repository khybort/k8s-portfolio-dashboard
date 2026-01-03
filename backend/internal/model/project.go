package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StringArray []string

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "[]", nil
	}
	return json.Marshal(a)
}

type Project struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null;index" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	GithubURL   string         `gorm:"type:varchar(500)" json:"github_url"`
	LiveURL     string         `gorm:"type:varchar(500)" json:"live_url"`
	Technologies StringArray   `gorm:"type:jsonb" json:"technologies"`
	Featured    bool           `gorm:"default:false;index:idx_projects_featured" json:"featured"`
	CreatedAt   time.Time      `gorm:"index:idx_projects_created_at" json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index:idx_projects_deleted_at" json:"-"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (p *Project) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Project) TableName() string {
	return "projects"
}

