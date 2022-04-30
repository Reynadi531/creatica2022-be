package services

import (
	"github.com/Reynadi531/creatica2022-be/internal/entities"
	"github.com/Reynadi531/creatica2022-be/internal/repository"
)

type CommentService interface {
	Save(comment entities.Comment) (entities.Comment, error)
	FindPostById(id string) (entities.Post, error)
}

type commentService struct {
	commentRepository repository.CommentRepositroy
}

func NewCommentService(r repository.CommentRepositroy) CommentService {
	return commentService{
		commentRepository: r,
	}
}

func (s commentService) Save(comment entities.Comment) (entities.Comment, error) {
	return s.commentRepository.Save(comment)
}

func (s commentService) FindPostById(id string) (entities.Post, error) {
	return s.commentRepository.FindPostById(id)
}
