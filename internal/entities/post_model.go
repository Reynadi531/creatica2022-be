package entities

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;not null"`
	Body      string    `gorm:"type:varchar(100);unique_index"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
