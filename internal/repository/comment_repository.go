package repository

import (
	"github.com/Reynadi531/creatica2022-be/internal/entities"
	"gorm.io/gorm"
)

type CommentRepositroy interface {
	Migrate() error
	Save(comment entities.Comment) (entities.Comment, error)
	FindPostById(id string) (entities.Post, error)
}

type commentRepository struct {
	DB *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepositroy {
	return commentRepository{
		DB: db,
	}
}

func (c commentRepository) Migrate() error {
	return c.DB.AutoMigrate(&entities.Comment{})
}

func (c commentRepository) Save(comment entities.Comment) (entities.Comment, error) {
	err := c.DB.Create(&comment).Error
	return comment, err
}

func (c commentRepository) FindPostById(id string) (entities.Post, error) {
	var post entities.Post
	err := c.DB.Where("id = ?", id).First(&post).Error
	return post, err
}
