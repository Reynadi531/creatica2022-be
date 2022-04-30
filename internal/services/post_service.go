package services

import (
	"github.com/Reynadi531/creatica2022-be/internal/entities"
	"github.com/Reynadi531/creatica2022-be/internal/repository"
	"github.com/Reynadi531/creatica2022-be/pkg/database"
)

type PostService interface {
	Save(post entities.Post) (entities.Post, error)
	List(pagination database.Pagination) (*database.Pagination, []*entities.Post, error)
	GetUserById(id string) (entities.User, error)
	GetPostById(id string) (entities.Post, error)
	GetCommentByPostId(id string) ([]*entities.Comment, error)
	CountCommentOnPost(id string) (int64, error)
}

type postService struct {
	postService repository.PostRepository
}

func NewPostService(r repository.PostRepository) PostService {
	return postService{
		postService: r,
	}
}

func (s postService) Save(post entities.Post) (entities.Post, error) {
	return s.postService.Save(post)
}

func (s postService) List(pagination database.Pagination) (*database.Pagination, []*entities.Post, error) {
	return s.postService.List(pagination)
}

func (s postService) GetUserById(id string) (entities.User, error) {
	return s.postService.GetUserById(id)
}

func (s postService) GetPostById(id string) (entities.Post, error) {
	return s.postService.GetPostById(id)
}

func (s postService) GetCommentByPostId(id string) ([]*entities.Comment, error) {
	return s.postService.FindCommentByPostID(id)
}

func (s postService) CountCommentOnPost(id string) (int64, error) {
	return s.postService.CountCommentOnPost(id)
}
