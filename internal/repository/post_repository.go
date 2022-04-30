package repository

import (
	"github.com/Reynadi531/creatica2022-be/internal/entities"
	"gorm.io/gorm"
)

type PostRepository interface {
	Migrate() error
	Save(entities.Post) (entities.Post, error)
}

type postRepository struct {
	DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return postRepository{
		DB: db,
	}
}

func (p postRepository) Migrate() error {
	return p.DB.AutoMigrate(&entities.Post{})
}

func (p postRepository) Save(post entities.Post) (entities.Post, error) {
	err := p.DB.Create(&post).Error
	return post, err
}
