package entities

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;not null"`
	Body      string    `gorm:"type:varchar(100);unique_index" json:"body" validate:"required"`
	Post      Post      `gorm:"foreignkey:PostID"`
	PostID    uuid.UUID `gorm:"type:uuid;not null" json:"post_id"`
	User      User      `gorm:"foreignKey:UserID"`
	UserID    uuid.UUID `gorm:"type;uuid;not null" json:"user_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
