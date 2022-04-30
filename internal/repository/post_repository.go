package repository

import (
	"github.com/Reynadi531/creatica2022-be/internal/entities"
	"github.com/Reynadi531/creatica2022-be/pkg/database"
	"gorm.io/gorm"
)

type PostRepository interface {
	Migrate() error
	Save(entities.Post) (entities.Post, error)
	List(database.Pagination) (*database.Pagination, []*entities.Post, error)
	GetUserById(id string) (entities.User, error)
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

func (p postRepository) List(pagination database.Pagination) (*database.Pagination, []*entities.Post, error) {
	var posts []*entities.Post
	p.DB.Scopes(database.Paginate(posts, &pagination, p.DB)).Find(&posts)
	return &pagination, posts, nil
}

func (p postRepository) GetUserById(id string) (entities.User, error) {
	var user entities.User
	err := p.DB.Where("id = ?", id).First(&user).Error
	return user, err
}
