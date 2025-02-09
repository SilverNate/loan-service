package loan

import "context"

//go:generate mockgen -destination=mocks/mock.go -package=mocks -source=interface.go

type IService interface {
	CreateLoan(ctx context.Context, loan *Loan) error
	GetLoanByID(ctx context.Context, id int64) (LoanDetail, error)
	ApproveLoan(ctx context.Context, approval ApprovalRequest) error
	InvestLoan(ctx context.Context, investment Investment) error
	DisburseLoan(ctx context.Context, disbursement Disbursement) error
}

type IRepository interface {
	CreateLoan(ctx context.Context, loan *Loan) error
	GetLoanByID(ctx context.Context, id int64) (*Loan, error)
	UpdateLoan(ctx context.Context, loan *Loan) error
	CreateInvestment(ctx context.Context, investment *Investment) error
	GetTotalInvestment(ctx context.Context, loanID int64) (float64, error)
	IsBorrowerExist(ctx context.Context, borrowerId int64) (isExist bool, err error)
	IsInvestorExist(ctx context.Context, investorId int64) (isExist bool, err error)
	GetInvestorByLoanId(ctx context.Context, loanId int64) (investorIDs []int, err error)
	GetInvestmentDetail(ctx context.Context, loanId int64) (detailInvestment []InvestmentDetail, err error)
}
