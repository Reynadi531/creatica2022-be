package services

import (
	"github.com/Reynadi531/creatica2022-be/internal/entities"
	"github.com/Reynadi531/creatica2022-be/internal/repository"
)

type AuthService interface {
	Save(entities.User) (entities.User, error)
	Update(entities.User) (entities.User, error)
	FindById(id string) (entities.User, error)
	FindByEmail(email string) (entities.User, error)
	FindByUsername(Username string) (entities.User, error)
	FindByRefreshToken(accessToken string) (entities.User, error)
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

func (s authService) FindByUsername(username string) (entities.User, error) {
	return s.authRepository.FindByUsername(username)
}

func (s authService) FindByRefreshToken(accessToken string) (entities.User, error) {
	return s.authRepository.FindByRefreshToken(accessToken)
}

func (s authService) Update(user entities.User) (entities.User, error) {
	return s.authRepository.Update(user)
}
