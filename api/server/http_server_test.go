package server

import (
	"TaskService/internal/repository"
	"TaskService/internal/repository/mocks"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockInterface := mocks.NewMockTaskServiceInterface(ctrl)
	mockRepo := NewTaskHandler(mockInterface)

	tests := []struct {
		name           string
		inputData      repository.Data
		expectedResult taskResult
		expectedError  bool
	}{
		{
			name:           "Success",
			inputData:      repository.Data{UserID: 1, Title: "Test", ID: 1, Description: "test"},
			expectedResult: taskResult{Message: "Task Created", Result: (*mongo.InsertOneResult)(nil)},
			expectedError:  false,
		},
		{
			name:           "Fail",
			inputData:      repository.Data{UserID: 1, Title: "Test", ID: 1, Description: "test"},
			expectedResult: taskResult{},
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		reqBody, _ := json.Marshal(tt.inputData)
		t.Run(tt.name, func(t *testing.T) {

			if tt.expectedError {
				mockInterface.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error some"))
			} else {
				mockInterface.EXPECT().CreateTask(gomock.Any(), gomock.Any())
			}

			req := httptest.NewRequest(http.MethodPost, "/createTask", bytes.NewReader(reqBody))
			rec := httptest.NewRecorder()

			mockRepo.createTask(rec, req)

			result := taskResult{}
			err := json.Unmarshal(rec.Body.Bytes(), &result)

			if tt.expectedError {
				assert.Error(t, err)
				assert.NotEqual(t, tt.expectedResult, &result, "expected and actual results do not match")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result, "expected and actual results do not match")
			}
		})
	}
}

func TestGetTaskByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskServiceInterface(ctrl)
	handler := NewTaskHandler(mockRepo)

	tests := []struct {
		name        string
		taskID      int64
		setupMock   func()
		expectedErr bool
	}{
		{
			name:   "Success",
			taskID: 1,
			setupMock: func() {
				expectedResult := repository.Data{ID: 1, Title: "Test Task", UserID: 1, Description: "Test Description"}
				mockRepo.EXPECT().GetTaskByID(gomock.Any(), gomock.Any()).Return(expectedResult, nil)
			},
			expectedErr: false,
		},
		{
			name:   "Failure - Not Found",
			taskID: 99,
			setupMock: func() {
				expectedResult := repository.Data{}
				mockRepo.EXPECT().GetTaskByID(gomock.Any(), gomock.Any()).Return(expectedResult, fmt.Errorf("error some"))
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/getTask?id="+strconv.FormatInt(tt.taskID, 10), nil)
			w := httptest.NewRecorder()
			if tt.setupMock != nil {
				tt.setupMock()
			}

			handler.getTaskByID(w, req)
			result := repository.Data{}
			err := json.Unmarshal(w.Body.Bytes(), &result)

			if tt.expectedErr {
				assert.NotEqual(t, tt.taskID, result.ID, "expected and actual results do not match")
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.taskID, result.ID, "expected and actual results do not match")
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAllTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockTaskServiceInterface(ctrl)
	handler := NewTaskHandler(mockRepo)
	tests := []struct {
		name        string
		setupMock   func()
		expectedErr bool
		lenResult   int
	}{
		{
			name: "Success",
			setupMock: func() {
				expectedResult := []repository.Data{
					{ID: 1, Title: "Test Task", UserID: 1, Description: "Test Description"},
					{ID: 2, Title: "Test Task", UserID: 1, Description: "Test Description"},
					{ID: 3, Title: "Test Task", UserID: 1, Description: "Test Description"},
					{ID: 4, Title: "Test Task", UserID: 1, Description: "Test Description"},
				}
				mockRepo.EXPECT().GetAllTasks(gomock.Any()).Return(expectedResult, nil)
			},
			expectedErr: false,
			lenResult:   4,
		},
		{
			name: "Failure - Not Found",
			setupMock: func() {
				var expectedResult []repository.Data
				mockRepo.EXPECT().GetAllTasks(gomock.Any()).Return(expectedResult, fmt.Errorf("error some"))
			},
			expectedErr: true,
			lenResult:   4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/getAllTasks", nil)
			w := httptest.NewRecorder()
			if tt.setupMock != nil {
				tt.setupMock()
			}
			handler.getAllTasks(w, req)
			var result []repository.Data
			err := json.Unmarshal(w.Body.Bytes(), &result)
			if tt.expectedErr {
				assert.NotEqual(t, tt.lenResult, len(result), "expected and actual results do not match")
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.lenResult, len(result), "expected and actual results do not match")
				assert.NoError(t, err)

			}
		})
	}

}

func TestStartServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		port        string
		expectError bool
	}{
		{"Valid Port", ":8081", false},
		{"Invalid Port", ":invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockTaskServiceInterface(ctrl)
			th := NewTaskHandler(mockRepo)
			_, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func() {
				err := th.StartServer()
				if tt.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
				cancel()
			}()

			time.Sleep(200 * time.Millisecond)

			cancel()
		})
	}
}
