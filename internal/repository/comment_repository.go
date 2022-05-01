package repository

import (
	"github.com/Reynadi531/creatica2022-be/internal/entities"
	"gorm.io/gorm"
)

type CommentRepositroy interface {
	Migrate() error
	Save(comment entities.Comment) (entities.Comment, error)
	FindPostById(id string) (entities.Post, error)
	FindCommentById(id string) (entities.Comment, error)
	FindReplyByCommentId(id string) ([]entities.Reply, error)
	SaveReply(reply entities.Reply) (entities.Reply, error)
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
	return c.DB.AutoMigrate(&entities.Comment{}, &entities.Reply{})
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

func (c commentRepository) FindCommentById(id string) (entities.Comment, error) {
	var comment entities.Comment
	err := c.DB.Where("id = ?", id).Preload("User").First(&comment).Error
	return comment, err
}

func (c commentRepository) FindReplyByCommentId(id string) ([]entities.Reply, error) {
	var reply []entities.Reply
	err := c.DB.Where("comment_id = ?", id).Preload("User").Find(&reply).Error
	return reply, err
}

func (c commentRepository) SaveReply(reply entities.Reply) (entities.Reply, error) {
	err := c.DB.Create(&reply).Error
	return reply, err
}
