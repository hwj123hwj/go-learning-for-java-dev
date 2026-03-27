package service

import (
	"errors"
	"user-api-advanced/model"
	"user-api-advanced/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	return s.repo.FindAll()
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	return s.repo.FindByEmail(email)
}

func (s *UserService) CreateUser(user *model.User) error {
	existing, _ := s.repo.FindByEmail(user.Email)
	if existing.ID != 0 {
		return errors.New("邮箱已存在")
	}
	return s.repo.Create(user)
}

func (s *UserService) UpdateUser(id uint, user *model.User) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if user.Name != "" {
		existing.Name = user.Name
	}
	if user.Email != "" {
		existing.Email = user.Email
	}
	if user.Age != 0 {
		existing.Age = user.Age
	}
	return s.repo.Update(existing)
}

func (s *UserService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
