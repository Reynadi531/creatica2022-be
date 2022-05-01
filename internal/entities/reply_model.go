package entities

import (
	"time"

	"github.com/google/uuid"
)

type Reply struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;not null"`
	Body      string    `gorm:"type:text;unique_index" json:"body" validate:"required"`
	User      User      `gorm:"foreignKey:UserID"`
	UserID    uuid.UUID `gorm:"type;uuid;not null" json:"user_id"`
	Comment   Comment   `gorm:"foreignKey:CommentID"`
	CommentID uuid.UUID `gorm:"type:uuid;not null" json:"comment_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
