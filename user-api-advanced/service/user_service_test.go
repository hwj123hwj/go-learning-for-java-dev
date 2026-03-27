package service

import (
	"testing"
	"user-api-advanced/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserService_GetAllUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	expectedUsers := []model.User{
		{ID: 1, Name: "Alice", Email: "alice@test.com"},
		{ID: 2, Name: "Bob", Email: "bob@test.com"},
	}

	mockRepo.On("FindAll").Return(expectedUsers, nil)

	users, err := userService.GetAllUsers()

	assert.NoError(t, err)
	assert.Equal(t, len(expectedUsers), len(users))
	assert.Equal(t, expectedUsers[0].Name, users[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	user := &model.User{Name: "Alice", Email: "alice@test.com", Age: 25}

	mockRepo.On("FindByEmail", "alice@test.com").Return(&model.User{}, nil)
	mockRepo.On("Create", user).Return(nil)

	err := userService.CreateUser(user)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_EmailExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	user := &model.User{Name: "Alice", Email: "alice@test.com", Age: 25}
	existingUser := &model.User{ID: 1, Name: "Alice", Email: "alice@test.com"}

	mockRepo.On("FindByEmail", "alice@test.com").Return(existingUser, nil)

	err := userService.CreateUser(user)

	assert.Error(t, err)
	assert.Equal(t, "邮箱已存在", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	mockRepo.On("Delete", uint(1)).Return(nil)

	err := userService.DeleteUser(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
