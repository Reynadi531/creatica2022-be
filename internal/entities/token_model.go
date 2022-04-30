package entities

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;not null"`
	Token     string    `gorm:"not null;unique"`
	UserId    uuid.UUID `gorm:"not null;unique"`
	User      User
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
