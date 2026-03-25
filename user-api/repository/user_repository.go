package repository

import (
	"user-api/config"
	"user-api/model"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindAll() ([]model.User, error) {
	var users []model.User
	err := config.DB.Find(&users).Error
	return users, err
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := config.DB.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) Create(user *model.User) error {
	return config.DB.Create(user).Error
}

func (r *UserRepository) Update(user *model.User) error {
	return config.DB.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return config.DB.Delete(&model.User{}, id).Error
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}
