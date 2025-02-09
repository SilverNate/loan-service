package loan_test

import (
	"bytes"
	"context"
	"errors"
	"loan-service/internal/loan"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLoanService struct {
	mock.Mock
}

func (m *MockLoanService) GetLoanByID(ctx context.Context, id int64) (loan.LoanDetail, error) {
	//TODO: create mock
	panic("implement me")
}

func (m *MockLoanService) ApproveLoan(ctx context.Context, approval loan.ApprovalRequest) error {
	args := m.Called(ctx, approval)
	return args.Error(0)
}

func (m *MockLoanService) InvestLoan(ctx context.Context, investment loan.Investment) error {
	args := m.Called(ctx, investment)
	return args.Error(0)
}

func (m *MockLoanService) DisburseLoan(ctx context.Context, disbursement loan.Disbursement) error {
	args := m.Called(ctx, disbursement)
	return args.Error(0)
}

func (m *MockLoanService) CreateLoan(ctx context.Context, loan *loan.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func TestCreateLoanHandler_Success(t *testing.T) {
	mockService := new(MockLoanService)
	handler := &loan.LoanHandler{LoanService: mockService}

	// Define the request body
	requestBody := `{
		"borrower_id": 1,
		"principal_amount": 1000000,
		"rate": 5.0,
		"loan_term": 12
	}`

	// Create a request with the body
	req, err := http.NewRequest(http.MethodPost, "/loans", bytes.NewBufferString(requestBody))
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Set up the mock service to return no error
	mockService.On("CreateLoan", mock.Anything, mock.AnythingOfType("*Loan")).Return(nil)

	// Call the handler
	handler.CreateLoan(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Body.String(), "create loan success")

	// Verify the mock was called
	mockService.AssertExpectations(t)
}

func TestCreateLoanHandler_InvalidRequestBody(t *testing.T) {
	// Create a mock LoanService
	mockService := new(MockLoanService)
	handler := &loan.LoanHandler{LoanService: mockService}

	// Define an invalid request body
	requestBody := `{invalid json}`

	// Create a request with the invalid body
	req, err := http.NewRequest(http.MethodPost, "/loans", bytes.NewBufferString(requestBody))
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.CreateLoan(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "invalid character")

	// Verify the mock was not called
	mockService.AssertNotCalled(t, "CreateLoan")
}

func TestCreateLoanHandler_ValidationError(t *testing.T) {
	// Create a mock LoanService
	mockService := new(MockLoanService)
	handler := &loan.LoanHandler{LoanService: mockService}

	// Define a request body with invalid data
	requestBody := `{
		"borrower_id": 0,
		"principal_amount": 0,
		"rate": 0,
		"loan_term": 0
	}`

	// Create a request with the body
	req, err := http.NewRequest(http.MethodPost, "/loans", bytes.NewBufferString(requestBody))
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.CreateLoan(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "borrower_id must be greater than 0")

	// Verify the mock was not called
	mockService.AssertNotCalled(t, "CreateLoan")
}

func TestCreateLoanHandler_ServiceError(t *testing.T) {
	// Create a mock LoanService
	mockService := new(MockLoanService)
	handler := &loan.LoanHandler{LoanService: mockService}

	// Define the request body
	requestBody := `{
		"borrower_id": 1,
		"principal_amount": 1000000,
		"rate": 5.0,
		"loan_term": 12
	}`

	// Create a request with the body
	req, err := http.NewRequest(http.MethodPost, "/loans", bytes.NewBufferString(requestBody))
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Set up the mock service to return an error
	mockService.On("CreateLoan", mock.Anything, mock.AnythingOfType("*Loan")).Return(errors.New("service error"))

	// Call the handler
	handler.CreateLoan(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "service error")

	// Verify the mock was called
	mockService.AssertExpectations(t)
}
