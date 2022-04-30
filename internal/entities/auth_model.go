package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;not null"`
	Email     string    `gorm:"type:varchar(100);unique_index" validate:"email"`
	Username  string    `gorm:"type:varchar(100);unique_index" validate:"required"`
	Password  string    `gorm:"type:varchar(100);unique_index" validate:"required"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
