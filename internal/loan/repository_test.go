package loan_test

import (
	"context"
	"errors"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"loan-service/internal/loan"
	"regexp"
	"testing"
)

func TestCreateLoan_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := loan.NewLoanRepository(db)

	loan := &loan.Loan{
		BorrowerID:                 1,
		PrincipalAmount:            1000000.0,
		Rate:                       5.0,
		TotalInterest:              50000.0,
		AgreementLetter:            nil,
		FieldValidatorPictureProof: nil,
		FieldValidatorID:           2,
		Status:                     "proposed",
		LoanTerm:                   12,
	}

	expectedQuery := regexp.QuoteMeta(`
    INSERT INTO loans (borrower_id, principal_amount, rate, total_interest, agreement_letter, field_validator_picture_proof, field_validator_id, status, loan_term, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    RETURNING id`)

	mock.ExpectQuery(expectedQuery).
		WithArgs(
			loan.BorrowerID,
			loan.PrincipalAmount,
			loan.Rate,
			loan.TotalInterest,
			loan.AgreementLetter,
			loan.FieldValidatorPictureProof,
			loan.FieldValidatorID,
			loan.Status,
			loan.LoanTerm,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = repo.CreateLoan(context.Background(), loan)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), loan.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateLoan_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := loan.NewLoanRepository(db)

	loan := &loan.Loan{
		BorrowerID:                 1,
		PrincipalAmount:            1000000.0,
		Rate:                       5.0,
		TotalInterest:              50000.0,
		AgreementLetter:            nil,
		FieldValidatorPictureProof: nil,
		FieldValidatorID:           2,
		Status:                     "proposed",
		LoanTerm:                   12,
	}

	expectedQuery := regexp.QuoteMeta(`
    INSERT INTO loans (borrower_id, principal_amount, rate, total_interest, agreement_letter, field_validator_picture_proof, field_validator_id, status, loan_term, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    RETURNING id`)
	mock.ExpectQuery(expectedQuery).
		WithArgs(
			loan.BorrowerID,
			loan.PrincipalAmount,
			loan.Rate,
			loan.TotalInterest,
			loan.AgreementLetter,
			loan.FieldValidatorPictureProof,
			loan.FieldValidatorID,
			loan.Status,
			loan.LoanTerm,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(errors.New("database error"))

	err = repo.CreateLoan(context.Background(), loan)

	assert.Error(t, err)
	assert.EqualError(t, err, "create loan: database error")

	assert.NoError(t, mock.ExpectationsWereMet())
}
