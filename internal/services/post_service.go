package services

import (
	"github.com/Reynadi531/creatica2022-be/internal/entities"
	"github.com/Reynadi531/creatica2022-be/internal/repository"
)

type PostService interface {
	Save(post entities.Post) (entities.Post, error)
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
