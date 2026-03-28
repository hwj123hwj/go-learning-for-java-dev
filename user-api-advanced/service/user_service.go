package service

import (
	"errors"
	"user-api-advanced/model"
	"user-api-advanced/repository"

	"gorm.io/gorm"
)

// UserService 定义业务层接口
type UserService interface {
	GetAllUsers() ([]model.User, error)
	GetUserByID(id uint) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(id uint, user *model.User) error
	DeleteUser(id uint) error
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService 依赖 UserRepository 接口
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAllUsers() ([]model.User, error) {
	return s.repo.FindAll()
}

func (s *userService) GetUserByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) GetUserByEmail(email string) (*model.User, error) {
	return s.repo.FindByEmail(email)
}

func (s *userService) CreateUser(user *model.User) error {
	_, err := s.repo.FindByEmail(user.Email)
	if err == nil {
		// 找到记录 → 邮箱已存在
		return errors.New("邮箱已存在")
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
	if user.Name != "" {
		existing.Name = user.Name
	}
	if user.Email != "" {
		existing.Email = user.Email
	}
	// Age 使用指针，nil 表示"不更新"，非 nil 则覆盖（包括更新为 0 的情况）
	if user.Age != nil {
		existing.Age = user.Age
	}
	return s.repo.Update(existing)
}

func (s *userService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
