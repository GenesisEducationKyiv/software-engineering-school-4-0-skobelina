package subscribers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/skobelina/currency_converter/domains/subscribers"
	utils "github.com/skobelina/currency_converter/utils/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Create(request *subscribers.SubscriberRequest) (*string, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockService) Search(filter *subscribers.SearchSubscribeRequest) (*subscribers.SearchSubscribeResponse, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*subscribers.SearchSubscribeResponse), args.Error(1)
}

func TestHandler_Create_Success(t *testing.T) {
	mockService := new(MockService)
	mockResponse := "Subscription created"
	mockService.On("Create", mock.Anything).Return(&mockResponse, nil)

	handler := subscribers.NewHandler(mockService)
	router := mux.NewRouter()
	handler.Register(router)

	requestBody, err := json.Marshal(subscribers.SubscriberRequest{})
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/subscribe", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	expectedResponse, err := json.Marshal(mockResponse)
	assert.NoError(t, err)
	assert.JSONEq(t, string(expectedResponse), rr.Body.String())
	mockService.AssertExpectations(t)
}

func TestHandler_Create_Conflict(t *testing.T) {
	mockService := new(MockService)
	mockService.On("Create", mock.Anything).Return(nil, utils.NewIsConflictError("conflict error"))

	handler := subscribers.NewHandler(mockService)
	router := mux.NewRouter()
	handler.Register(router)

	requestBody, err := json.Marshal(subscribers.SubscriberRequest{})
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/subscribe", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
	expectedResponse := `{"message": "conflict error", "type": "Conflict error"}`
	assert.JSONEq(t, expectedResponse, rr.Body.String())
	mockService.AssertExpectations(t)
}

func TestHandler_Search_Success(t *testing.T) {
	mockService := new(MockService)
	mockResponse := &subscribers.SearchSubscribeResponse{}
	mockService.On("Search", mock.Anything).Return(mockResponse, nil)

	handler := subscribers.NewHandler(mockService)
	router := mux.NewRouter()
	handler.Register(router)

	req, err := http.NewRequest("GET", "/api/subscribe", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	expectedResponse, err := json.Marshal(mockResponse)
	assert.NoError(t, err)
	assert.JSONEq(t, string(expectedResponse), rr.Body.String())
	mockService.AssertExpectations(t)
}

func TestHandler_Search_BadRequest(t *testing.T) {
	mockService := new(MockService)
	mockService.On("Search", mock.Anything).Return(nil, utils.NewBadRequestError("bad request"))

	handler := subscribers.NewHandler(mockService)
	router := mux.NewRouter()
	handler.Register(router)

	req, err := http.NewRequest("GET", "/api/subscribe", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expectedResponse := `{"message": "bad request", "type": "Bad request"}`
	assert.JSONEq(t, expectedResponse, rr.Body.String())
	mockService.AssertExpectations(t)
}
