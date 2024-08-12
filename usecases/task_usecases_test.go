// usecases/task_usecases_test.go
package usecases

import (
	"clean-architecture/domain"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTaskRepository is a mock implementation of the TaskRepository interface
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) GetTasks() ([]domain.Task, error) {
	args := m.Called()
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskRepository) GetTaskByID(id string) (domain.Task, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskRepository) CreateTask(task domain.Task) (domain.Task, error) {
	args := m.Called(task)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskRepository) UpdateTask(id string, task domain.Task) (domain.Task, error) {
	args := m.Called(id, task)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskRepository) DeleteTask(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Test for GetTasks method
func TestTaskUsecase_GetTasks(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskUsecase := NewTaskUsecase(mockRepo)

	expectedTasks := []domain.Task{
		{Title: "Task 1", Description: "Description 1"},
		{Title: "Task 2", Description: "Description 2"},
	}

	mockRepo.On("GetTasks").Return(expectedTasks, nil)

	tasks, err := taskUsecase.GetTasks()

	assert.NoError(t, err)
	assert.Equal(t, expectedTasks, tasks)
	mockRepo.AssertExpectations(t)
}

// Test for GetTaskByID method
func TestTaskUsecase_GetTaskByID(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskUsecase := NewTaskUsecase(mockRepo)

	expectedTask := domain.Task{Title: "Task 1", Description: "Description 1"}

	mockRepo.On("GetTaskByID", "1").Return(expectedTask, nil)

	task, err := taskUsecase.GetTaskByID("1")

	assert.NoError(t, err)
	assert.Equal(t, expectedTask, task)
	mockRepo.AssertExpectations(t)
}

// Test for CreateTask method
func TestTaskUsecase_CreateTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskUsecase := NewTaskUsecase(mockRepo)

	newTask := domain.Task{Title: "New Task", Description: "New Description"}
	createdTask := domain.Task{Title: "New Task", Description: "New Description"}

	mockRepo.On("CreateTask", newTask).Return(createdTask, nil)

	task, err := taskUsecase.CreateTask(newTask)

	assert.NoError(t, err)
	assert.Equal(t, createdTask, task)
	mockRepo.AssertExpectations(t)
}

// Test for UpdateTask method
func TestTaskUsecase_UpdateTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskUsecase := NewTaskUsecase(mockRepo)

	updatedTask := domain.Task{Title: "Updated Task", Description: "Updated Description"}

	mockRepo.On("UpdateTask", "1", updatedTask).Return(updatedTask, nil)

	task, err := taskUsecase.UpdateTask("1", updatedTask)

	assert.NoError(t, err)
	assert.Equal(t, updatedTask, task)
	mockRepo.AssertExpectations(t)
}

// Test for DeleteTask method
func TestTaskUsecase_DeleteTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskUsecase := NewTaskUsecase(mockRepo)

	mockRepo.On("DeleteTask", "1").Return(nil)

	err := taskUsecase.DeleteTask("1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Test for GetTaskByID with error
func TestTaskUsecase_GetTaskByID_Error(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskUsecase := NewTaskUsecase(mockRepo)

	mockRepo.On("GetTaskByID", "1").Return(domain.Task{}, errors.New("task not found"))

	_, err := taskUsecase.GetTaskByID("1")

	assert.Error(t, err)
	assert.Equal(t, "task not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// Test for CreateTask with error
func TestTaskUsecase_CreateTask_Error(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskUsecase := NewTaskUsecase(mockRepo)

	newTask := domain.Task{Title: "New Task", Description: "New Description"}

	mockRepo.On("CreateTask", newTask).Return(domain.Task{}, errors.New("failed to create task"))

	_, err := taskUsecase.CreateTask(newTask)

	assert.Error(t, err)
	assert.Equal(t, "failed to create task", err.Error())
	mockRepo.AssertExpectations(t)
}
