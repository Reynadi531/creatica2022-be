package repository

import (
	"github.com/Reynadi531/creatica2022-be/internal/entities"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Migrate() error
	Save(entities.User) (entities.User, error)
	Update(entities.User) (entities.User, error)
	FindById(id string) (entities.User, error)
	FindByEmail(email string) (entities.User, error)
	FindByUsername(Username string) (entities.User, error)
	FindByRefreshToken(accessToken string) (entities.User, error)
}

type authRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return authRepository{
		DB: db,
	}
}

func (u authRepository) Migrate() error {
	return u.DB.AutoMigrate(&entities.User{})
}

func (u authRepository) Save(user entities.User) (entities.User, error) {
	err := u.DB.Create(&user).Error
	return user, err
}

func (u authRepository) FindById(id string) (entities.User, error) {
	var user entities.User
	err := u.DB.Where("id = ?", id).First(&user).Error
	return user, err
}

func (u authRepository) FindByEmail(email string) (entities.User, error) {
	var user entities.User
	err := u.DB.Where("email = ?", email).First(&user).Error
	return user, err
}

func (u authRepository) FindByUsername(Username string) (entities.User, error) {
	var user entities.User
	err := u.DB.Where("username = ?", Username).First(&user).Error
	return user, err
}

func (u authRepository) FindByRefreshToken(accessToken string) (entities.User, error) {
	var user entities.User
	err := u.DB.Where("refresh_token = ?", accessToken).First(&user).Error
	return user, err
}

func (u authRepository) Update(user entities.User) (entities.User, error) {
	err := u.DB.Save(&user).Error
	return user, err
}
