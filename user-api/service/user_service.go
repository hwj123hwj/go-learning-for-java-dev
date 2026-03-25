package service

import (
	"errors"
	"user-api/model"
	"user-api/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	return s.repo.FindAll()
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) CreateUser(user *model.User) error {
	existing, _ := s.repo.FindByEmail(user.Email)
	if existing.ID != 0 {
		return errors.New("email already exists")
	}
	return s.repo.Create(user)
}

func (s *UserService) UpdateUser(id uint, user *model.User) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	existing.Name = user.Name
	existing.Email = user.Email
	existing.Age = user.Age
	return s.repo.Update(existing)
}

func (s *UserService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
