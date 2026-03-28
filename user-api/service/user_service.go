package service

import (
	"errors"
	"user-api/model"
	"user-api/repository"

	"gorm.io/gorm"
)

// UserService 定义业务层接口
type UserService interface {
	GetAllUsers() ([]model.User, error)
	GetUserByID(id uint) (*model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(id uint, user *model.User) error
	DeleteUser(id uint) error
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService 依赖 UserRepository 接口，而非具体结构体
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAllUsers() ([]model.User, error) {
	return s.repo.FindAll()
}

func (s *userService) GetUserByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) CreateUser(user *model.User) error {
	_, err := s.repo.FindByEmail(user.Email)
	if err == nil {
		// 找到了记录，说明邮箱已被注册
		return errors.New("email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 真正的数据库错误
		return err
	}
	return s.repo.Create(user)
}

func (s *userService) UpdateUser(id uint, user *model.User) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	// 只更新非零值字段，避免客户端漏传字段时把数据清空
	if user.Name != "" {
		existing.Name = user.Name
	}
	if user.Email != "" {
		existing.Email = user.Email
	}
	if user.Age != nil {
		existing.Age = user.Age
	}
	return s.repo.Update(existing)
}

func (s *userService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
