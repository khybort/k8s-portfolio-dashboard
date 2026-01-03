package model

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Article struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null;index" json:"title"`
	Slug        string         `gorm:"type:varchar(255);uniqueIndex:idx_articles_slug;not null" json:"slug"`
	Excerpt     string         `gorm:"type:text" json:"excerpt"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	AuthorID    uuid.UUID      `gorm:"type:uuid;not null;index:idx_articles_author_id" json:"author_id"`
	Published   bool           `gorm:"default:false;index:idx_articles_published" json:"published"`
	PublishedAt *time.Time     `gorm:"index:idx_articles_published_at" json:"published_at,omitempty"`
	CreatedAt   time.Time      `gorm:"index:idx_articles_created_at" json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index:idx_articles_deleted_at" json:"-"`
}

func (a *Article) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

func (a *Article) BeforeUpdate(tx *gorm.DB) error {
	a.UpdatedAt = time.Now()
	return nil
}

func (a *Article) TableName() string {
	return "articles"
}

