package services

import (
	"github.com/Reynadi531/creatica2022-be/internal/entities"
	"github.com/Reynadi531/creatica2022-be/internal/repository"
)

type AuthService interface {
	Save(entities.User) (entities.User, error)
	FindById(id string) (entities.User, error)
	FindByEmail(email string) (entities.User, error)
}

type authService struct {
	authRepository repository.AuthRepository
}

func NewAuthService(r repository.AuthRepository) AuthService {
	return authService{
		authRepository: r,
	}
}

func (s authService) FindByEmail(email string) (entities.User, error) {
	return s.authRepository.FindByEmail(email)
}

func (s authService) FindById(id string) (entities.User, error) {
	return s.authRepository.FindById(id)
}

func (s authService) Save(user entities.User) (entities.User, error) {
	return s.authRepository.Save(user)
}
