package rates

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/skobelina/currency_converter/internal/rates"
	utils "github.com/skobelina/currency_converter/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Get() (*float64, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*float64), args.Error(1)
}

func TestHandler_GetRate_Success(t *testing.T) {
	mockService := new(MockService)
	mockResponse := 27.32
	mockService.On("Get").Return(&mockResponse, nil)

	handler := rates.NewHandler(mockService)
	router := mux.NewRouter()
	handler.Register(router)

	req, err := http.NewRequest("GET", "/api/rate", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	expectedResponse, err := json.Marshal(mockResponse)
	assert.NoError(t, err)
	assert.JSONEq(t, string(expectedResponse), rr.Body.String())
	mockService.AssertExpectations(t)
}

func TestHandler_GetRate_InternalServerError(t *testing.T) {
	mockService := new(MockService)
	mockService.On("Get").Return(nil, utils.NewInternalServerError("internal server error"))

	handler := rates.NewHandler(mockService)
	router := mux.NewRouter()
	handler.Register(router)

	req, err := http.NewRequest("GET", "/api/rate", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expectedResponse := `{"message": "internal server error", "type": "Internal server error"}`
	assert.JSONEq(t, expectedResponse, rr.Body.String())
	mockService.AssertExpectations(t)
}

func TestHandler_GetRate_BadRequest(t *testing.T) {
	mockService := new(MockService)
	mockService.On("Get").Return(nil, utils.NewBadRequestError("bad request"))

	handler := rates.NewHandler(mockService)
	router := mux.NewRouter()
	handler.Register(router)

	req, err := http.NewRequest("GET", "/api/rate", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expectedResponse := `{"message": "bad request", "type": "Bad request"}`
	assert.JSONEq(t, expectedResponse, rr.Body.String())
	mockService.AssertExpectations(t)
}
