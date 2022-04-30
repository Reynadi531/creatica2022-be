package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;not null"`
	Username     string    `gorm:"type:varchar(100);unique_index" validate:"required"`
	Password     string    `gorm:"type:varchar(100);unique_index" validate:"required"`
	RefreshToken string    `gorm:"type:varchar(100);unique_index"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
