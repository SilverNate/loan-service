package loan

import (
	"database/sql"
	"time"
)

type Loan struct {
	ID                         int64        `json:"id"`
	BorrowerID                 int64        `json:"borrower_id"`
	PrincipalAmount            float64      `json:"principal_amount"`
	Rate                       float64      `json:"rate"`
	TotalInterest              float64      `json:"total_interest"`
	AgreementLetter            *string      `json:"agreement_letter"`
	LoanTerm                   int          `json:"loan_term"`
	FieldValidatorPictureProof *string      `json:"field_validator_picture_proof"`
	FieldValidatorID           int          `json:"field_validator_id"`
	FieldOfficerId             int          `json:"field_officer_id"`
	ApprovalDate               sql.NullTime `json:"approval_date"`
	DisbursementDate           time.Time    `json:"disbursement_date"`
	Status                     string       `json:"status"` // proposed, approved, invested, disbursed
	CreatedAt                  time.Time    `json:"created_at"`
	UpdatedAt                  time.Time    `json:"updated_at"`
}

// not used we use hardcode for this system right now
type Borrower struct {
	ID              int64     `json:"id"`
	AgreementLetter string    `json:"agreement_letter"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// not used we use hardcode for this system right now
type Investor struct {
	ID              int64     `json:"id"`
	AgreementLetter string    `json:"agreement_letter"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Disbursement struct {
	LoanID          int64  `json:"loan_id"`
	AgreementLetter string `json:"agreement_letter"`
	FieldOfficerID  int    `json:"field_officer_id"`
}

type Investment struct {
	ID         int64     `json:"id"`
	InvestorID int64     `json:"investor_id"`
	LoanID     int64     `json:"loan_id"`
	Amount     float64   `json:"amount"`
	ROI        float64   `json:"roi"`
	CreatedAt  time.Time `json:"created_at"`
}
