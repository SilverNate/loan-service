package loan

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type LoanRepository struct {
	db *sql.DB
}

func NewLoanRepository(db *sql.DB) *LoanRepository {
	return &LoanRepository{db: db}
}

func (r *LoanRepository) CreateLoan(ctx context.Context, loan *Loan) error {
	query := `
		INSERT INTO loans (borrower_id, principal_amount, rate, total_interest,  agreement_letter, field_validator_picture_proof, field_validator_id, status, loan_term, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	now := time.Now()
	err := r.db.QueryRowContext(ctx, query, loan.BorrowerID, loan.PrincipalAmount, loan.Rate, loan.TotalInterest, loan.AgreementLetter, loan.FieldValidatorPictureProof, loan.FieldValidatorID, loan.Status, loan.LoanTerm, now, now).Scan(&loan.ID)
	if err != nil {
		logrus.Errorf("error create loan: %v", err.Error())
		return fmt.Errorf("create loan: %v", err)
	}

	return nil
}

func (r *LoanRepository) GetLoanByID(ctx context.Context, id int64) (*Loan, error) {
	query := `
		SELECT id, borrower_id, principal_amount, rate, agreement_letter, status, created_at, updated_at, total_interest, field_validator_picture_proof, field_validator_id, approval_date
		FROM loans WHERE id = $1`
	var loan Loan
	err := r.db.QueryRowContext(ctx, query, id).Scan(&loan.ID, &loan.BorrowerID, &loan.PrincipalAmount, &loan.Rate, &loan.AgreementLetter, &loan.Status, &loan.CreatedAt, &loan.UpdatedAt, &loan.TotalInterest, &loan.FieldValidatorPictureProof, &loan.FieldValidatorID, &loan.ApprovalDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("loan not found")
		}
		logrus.Errorf("error get loan by id: %v", err.Error())
		return nil, err
	}
	return &loan, nil
}

func (r *LoanRepository) UpdateLoan(ctx context.Context, loan *Loan) error {
	loan.UpdatedAt = time.Now()
	query := `
		UPDATE loans
		SET borrower_id=$1, principal_amount=$2, rate=$3, agreement_letter=$4, status=$5, updated_at=$6, approval_date = $7, field_officer_id = $8, field_validator_picture_proof = $9, field_validator_id=$10, disbursment_date=$11
		WHERE id=$12`
	_, err := r.db.ExecContext(ctx, query, loan.BorrowerID, loan.PrincipalAmount, loan.Rate, loan.AgreementLetter, loan.Status, loan.UpdatedAt, loan.ApprovalDate, loan.FieldOfficerId, loan.FieldValidatorPictureProof, loan.FieldValidatorID, loan.DisbursementDate, loan.ID)
	if err != nil {
		logrus.Errorf("error update loan: %v", err.Error())
		return err
	}

	return nil
}

func (r *LoanRepository) CreateInvestment(ctx context.Context, inv *Investment) error {
	query := `
		INSERT INTO investments (loan_id, investor_id, amount, roi, total_gain, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	now := time.Now()
	err := r.db.QueryRowContext(ctx, query, inv.LoanID, inv.InvestorID, inv.Amount, inv.ROI, inv.TotalGain, now).Scan(&inv.ID)
	if err != nil {
		logrus.Errorf("error create investment: %v", err.Error())
		return err
	}

	return nil
}

func (r *LoanRepository) GetTotalInvestment(ctx context.Context, loanID int64) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM investments WHERE loan_id = $1`
	var sum float64
	err := r.db.QueryRowContext(ctx, query, loanID).Scan(&sum)
	if err != nil {
		logrus.Errorf("error get total invesment: %v", err.Error())
		return sum, err
	}
	return sum, err
}

func (r *LoanRepository) IsBorrowerExist(ctx context.Context, borrowerId int64) (isExist bool, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM borrowers WHERE id = $1)`

	err = r.db.QueryRowContext(ctx, query, borrowerId).Scan(&isExist)
	if err != nil {
		logrus.Errorf("error checking borrower is exist: %v", err.Error())
		return isExist, err
	}

	return isExist, nil
}

func (r *LoanRepository) IsInvestorExist(ctx context.Context, investorId int64) (isExist bool, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM investors WHERE id = $1)`

	err = r.db.QueryRowContext(ctx, query, investorId).Scan(&isExist)
	if err != nil {
		logrus.Errorf("error checking investor is exist: %v", err.Error())
		return isExist, err
	}

	return isExist, nil
}

func (r *LoanRepository) GetInvestorByLoanId(ctx context.Context, loanId int64) (investorIDs []int, err error) {
	query := `
		SELECT investor_id
		FROM investments
		WHERE loan_id = $1`

	rows, err := r.db.QueryContext(ctx, query, loanId)
	if err != nil {
		logrus.Errorf("error get investor by loan_id: %v", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var investorID int
		if err := rows.Scan(&investorID); err != nil {
			logrus.Errorf("error scan get investor by loan_id: %v", err.Error())
			return nil, err
		}
		investorIDs = append(investorIDs, investorID)
	}

	if err := rows.Err(); err != nil {
		logrus.Errorf("rows error investor by load_id: %v", err.Error())
		return nil, err
	}

	return investorIDs, nil
}

func (r *LoanRepository) GetInvestmentDetail(ctx context.Context, loanId int64) (detailInvestment []InvestmentDetail, err error) {
	query := `
		SELECT investor_id, SUM(amount) AS total_amount, sum(roi) as total_roi
		FROM investments
		WHERE loan_id = $1
		GROUP BY investor_id, loan_id`
	rows, err := r.db.QueryContext(ctx, query, loanId)
	if err != nil {
		logrus.Errorf("error query get investment detail: %v", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var summary InvestmentDetail
		if err := rows.Scan(&summary.InvestorId, &summary.TotalAmount, &summary.TotalROI); err != nil {
			logrus.Errorf("error scan get investment detail: %v", err.Error())
			return nil, err
		}
		detailInvestment = append(detailInvestment, summary)
	}

	if err := rows.Err(); err != nil {
		logrus.Errorf("rows error get investment detail: %v", err.Error())
		return nil, err
	}

	return
}
