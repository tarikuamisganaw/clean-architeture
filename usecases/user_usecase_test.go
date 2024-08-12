// usecases/user_usecases_test.go
package usecases

import (
	"clean-architecture/domain"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Register(user domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (domain.User, error) {
	args := m.Called(username)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUsers() ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

// MockPasswordService is a mock implementation of the PasswordService interface
type MockPasswordService struct {
	mock.Mock
}

func (m *MockPasswordService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordService) CheckPasswordHash(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

// MockJWTService is a mock implementation of the JWTService interface
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateJWT(username, role string) (string, error) {
	args := m.Called(username, role)
	return args.String(0), args.Error(1)
}

// Test for Register method
func TestUserUsecase_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockPasswordSvc := new(MockPasswordService)
	mockJWTService := new(MockJWTService)
	userUsecase := NewUserUsecase(mockRepo, mockPasswordSvc, mockJWTService)

	user := domain.User{Username: "testuser", Password: "password"}
	hashedPassword := "hashedpassword"

	mockPasswordSvc.On("HashPassword", user.Password).Return(hashedPassword, nil)
	mockRepo.On("Register", mock.AnythingOfType("domain.User")).Return(domain.User{Username: user.Username, Password: hashedPassword}, nil)

	createdUser, err := userUsecase.Register(user)

	assert.NoError(t, err)
	assert.Equal(t, hashedPassword, createdUser.Password)
	mockPasswordSvc.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// Test for Login method
func TestUserUsecase_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockPasswordSvc := new(MockPasswordService)
	mockJWTService := new(MockJWTService)
	userUsecase := NewUserUsecase(mockRepo, mockPasswordSvc, mockJWTService)

	user := domain.User{Username: "testuser", Password: "hashedpassword", Role: "user"}
	token := "jwt_token"

	mockRepo.On("FindByUsername", user.Username).Return(user, nil)
	mockPasswordSvc.On("CheckPasswordHash", user.Password, "password").Return(nil)
	mockJWTService.On("GenerateJWT", user.Username, user.Role).Return(token, nil)

	jwtToken, err := userUsecase.Login(user.Username, "password")

	assert.NoError(t, err)
	assert.Equal(t, token, jwtToken)
	mockRepo.AssertExpectations(t)
	mockPasswordSvc.AssertExpectations(t)
	mockJWTService.AssertExpectations(t)
}

// Test for GetUsers method
func TestUserUsecase_GetUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockPasswordSvc := new(MockPasswordService)
	mockJWTService := new(MockJWTService)
	userUsecase := NewUserUsecase(mockRepo, mockPasswordSvc, mockJWTService)

	expectedUsers := []domain.User{
		{Username: "user1", Password: "password1"},
		{Username: "user2", Password: "password2"},
	}

	mockRepo.On("GetUsers").Return(expectedUsers, nil)

	users, err := userUsecase.GetUsers()

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockRepo.AssertExpectations(t)
}

// Test for Login with invalid username or password
func TestUserUsecase_Login_InvalidCredentials(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockPasswordSvc := new(MockPasswordService)
	mockJWTService := new(MockJWTService)
	userUsecase := NewUserUsecase(mockRepo, mockPasswordSvc, mockJWTService)

	mockRepo.On("FindByUsername", "invaliduser").Return(domain.User{}, errors.New("invalid username or password"))

	_, err := userUsecase.Login("invaliduser", "wrongpassword")

	assert.Error(t, err)
	assert.Equal(t, "invalid username or password", err.Error())
	mockRepo.AssertExpectations(t)
}

// Test for Register with error hashing password
func TestUserUsecase_Register_ErrorHashingPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockPasswordSvc := new(MockPasswordService)
	mockJWTService := new(MockJWTService)
	userUsecase := NewUserUsecase(mockRepo, mockPasswordSvc, mockJWTService)

	user := domain.User{Username: "testuser", Password: "password"}

	mockPasswordSvc.On("HashPassword", user.Password).Return("", errors.New("hashing error"))

	_, err := userUsecase.Register(user)

	assert.Error(t, err)
	assert.Equal(t, "hashing error", err.Error())
	mockPasswordSvc.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Register", user)
}
