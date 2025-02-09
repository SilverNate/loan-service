package loan

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type CreateLoanRequest struct {
	BorrowerId      int64   `json:"borrower_id" validate:"required"`
	PrincipalAmount float64 `json:"principal_amount" validate:"required,gt=0"`
	Rate            float64 `json:"rate" validate:"required,gt=0"`
	LoanTerm        int     `json:"loan_term" validate:"required,min=1"`
}

func (request *CreateLoanRequest) ValidateCreateLoan() (err error) {
	validate := validator.New()
	err = validate.Struct(request)
	return
}

type ApprovalRequest struct {
	LoanId                int64  `json:"loan_id" validate:"required"`
	FieldValidatorPicTure string `json:"field_validator_picture" validate:"required"`
	FieldValidatorID      int    `json:"field_validator_id" validate:"required"`
}

func (request *ApprovalRequest) ValidateApprovalRequest() (err error) {
	validate := validator.New()
	err = validate.Struct(request)
	return
}

type InvestmentRequest struct {
	LoanId     int64   `json:"loan_id" validate:"required"`
	InvestorId int64   `json:"investor_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
}

func (request *InvestmentRequest) ValidateInvestmentRequest() (err error) {
	validate := validator.New()
	err = validate.Struct(request)
	return
}

type LoanDetail struct {
	ID                         int64              `json:"id"`
	BorrowerID                 int64              `json:"borrower_id"`
	PrincipalAmount            float64            `json:"principal_amount"`
	Rate                       float64            `json:"rate"`
	TotalInterest              float64            `json:"total_interest"`
	AgreementLetter            *string            `json:"agreement_letter"`
	LoanTerm                   int                `json:"loan_term"`
	FieldValidatorPictureProof *string            `json:"field_validator_picture_proof"`
	FieldValidatorID           int                `json:"field_validator_id"`
	FieldOfficerId             int                `json:"field_officer_id"`
	ApprovalDate               time.Time          `json:"approval_date"`
	DisbursementDate           time.Time          `json:"disbursement_date"`
	Status                     string             `json:"status"`
	CreatedAt                  time.Time          `json:"created_at"`
	UpdatedAt                  time.Time          `json:"updated_at"`
	InvestmentDetail           []InvestmentDetail `json:"investment_detail"`
}

type InvestmentDetail struct {
	InvestorId  int64   `json:"investor_id"`
	TotalAmount float64 `json:"total_amount"`
	TotalROI    float64 `json:"total_roi"`
}
