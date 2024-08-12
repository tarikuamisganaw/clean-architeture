// controllers/controller_test.go
package controllers

import (
	"clean-architecture/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockTaskUsecase is a mock implementation of the TaskUsecase interface
type MockTaskUsecase struct {
	mock.Mock
}

func (m *MockTaskUsecase) GetTasks() ([]domain.Task, error) {
	args := m.Called()
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) GetTaskByID(id string) (domain.Task, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) CreateTask(task domain.Task) (domain.Task, error) {
	args := m.Called(task)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) UpdateTask(id string, task domain.Task) (domain.Task, error) {
	args := m.Called(id, task)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) DeleteTask(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockUserUsecase is a mock implementation of the UserUsecase interface
type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Register(user domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserUsecase) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *MockUserUsecase) GetUsers() ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

// Test for GetTasks method
func TestTaskController_GetTasks(t *testing.T) {
	mockTaskUsecase := new(MockTaskUsecase)
	taskID1 := primitive.NewObjectID()
	taskID2 := primitive.NewObjectID()
	tasks := []domain.Task{

		{ID: taskID1, Title: "Task 1"},
		{ID: taskID2, Title: "Task 2"},
	}
	mockTaskUsecase.On("GetTasks").Return(tasks, nil)

	taskController := NewTaskController(mockTaskUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	taskController.GetTasks(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `[{"ID":"1","Title":"Task 1"},{"ID":"2","Title":"Task 2"}]`, w.Body.String())
	mockTaskUsecase.AssertExpectations(t)
}

// Test for GetTaskByID method
func TestTaskController_GetTaskByID(t *testing.T) {
	mockTaskUsecase := new(MockTaskUsecase)
	taskID1 := primitive.NewObjectID()
	task := domain.Task{ID: taskID1, Title: "Task 1"}
	mockTaskUsecase.On("GetTaskByID", "1").Return(task, nil)

	taskController := NewTaskController(mockTaskUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	taskController.GetTaskByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"ID":"1","Title":"Task 1"}`, w.Body.String())
	mockTaskUsecase.AssertExpectations(t)
}

// Test for CreateTask method
func TestTaskController_CreateTask(t *testing.T) {
	mockTaskUsecase := new(MockTaskUsecase)
	task := domain.Task{Title: "Task 1"}
	taskID1 := primitive.NewObjectID()
	createdTask := domain.Task{ID: taskID1, Title: "Task 1"}
	mockTaskUsecase.On("CreateTask", task).Return(createdTask, nil)

	taskController := NewTaskController(mockTaskUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/tasks", nil)
	c.Set("Content-Type", "application/json")
	c.Set("Body", task)

	taskController.CreateTask(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{"ID":"1","Title":"Task 1"}`, w.Body.String())
	mockTaskUsecase.AssertExpectations(t)
}

// Test for UpdateTask method
func TestTaskController_UpdateTask(t *testing.T) {
	mockTaskUsecase := new(MockTaskUsecase)
	task := domain.Task{Title: "Updated Task"}
	taskID1 := primitive.NewObjectID()
	updatedTask := domain.Task{ID: taskID1, Title: "Updated Task"}
	mockTaskUsecase.On("UpdateTask", taskID1, task).Return(updatedTask, nil)

	taskController := NewTaskController(mockTaskUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("PUT", "/tasks/1", nil)
	c.Set("Content-Type", "application/json")
	c.Set("Body", task)

	taskController.UpdateTask(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"ID":"1","Title":"Updated Task"}`, w.Body.String())
	mockTaskUsecase.AssertExpectations(t)
}

// Test for DeleteTask method
func TestTaskController_DeleteTask(t *testing.T) {
	mockTaskUsecase := new(MockTaskUsecase)
	mockTaskUsecase.On("DeleteTask", "1").Return(nil)
	taskID1 := primitive.NewObjectID()

	taskController := NewTaskController(mockTaskUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	taskController.DeleteTask(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message":"Task deleted successfully"}`, w.Body.String())
	mockTaskUsecase.AssertExpectations(t)
}

// Test for Register method
func TestUserController_Register(t *testing.T) {
	mockUserUsecase := new(MockUserUsecase)
	user := domain.User{Username: "testuser", Password: "password123"}
	taskID1 := primitive.NewObjectID()
	registeredUser := domain.User{ID: taskID1, Username: "testuser", Password: "hashedpassword123"}
	mockUserUsecase.On("Register", user).Return(registeredUser, nil)

	userController := NewUserController(mockUserUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/register", nil)
	c.Set("Content-Type", "application/json")
	c.Set("Body", user)

	userController.Register(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{"ID":"1","Username":"testuser","Password":"hashedpassword123"}`, w.Body.String())
	mockUserUsecase.AssertExpectations(t)
}

// Test for Login method
func TestUserController_Login(t *testing.T) {
	mockUserUsecase := new(MockUserUsecase)
	user := domain.User{Username: "testuser", Password: "password123"}
	token := "token123"
	mockUserUsecase.On("Login", user.Username, user.Password).Return(token, nil)

	userController := NewUserController(mockUserUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/login", nil)
	c.Set("Content-Type", "application/json")
	c.Set("Body", user)

	userController.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"token":"token123"}`, w.Body.String())
	mockUserUsecase.AssertExpectations(t)
}

// Test for GetUsers method
func TestUserController_GetUsers(t *testing.T) {
	mockUserUsecase := new(MockUserUsecase)
	taskID1 := primitive.NewObjectID()
	taskID2 := primitive.NewObjectID()

	users := []domain.User{
		{ID: taskID1, Username: "user1"},
		{ID: taskID2, Username: "user2"},
	}
	mockUserUsecase.On("GetUsers").Return(users, nil)

	userController := NewUserController(mockUserUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	userController.GetUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `[{"ID":"1","Username":"user1"},{"ID":"2","Username":"user2"}]`, w.Body.String())
	mockUserUsecase.AssertExpectations(t)
}
