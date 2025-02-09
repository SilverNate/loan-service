package loan_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	mockEvent "loan-service/internal/events/mocks"
	"loan-service/internal/loan"
	"loan-service/internal/loan/mocks"
	mockEmail "loan-service/internal/pkg/emailService/mocks"
	"testing"
)

type LoanSuite struct {
	suite.Suite
	ctx       context.Context
	svc       *loan.LoanService
	mockRepo  *mocks.MockIRepository
	mockEvent *mockEvent.MockEventPublisher
	mockEmail *mockEmail.MockEmailService
}

func (s *LoanSuite) SetupSuite() {
	mockCtrl := gomock.NewController(s.T())
	defer mockCtrl.Finish()
	s.ctx = context.Background()
	s.mockRepo = mocks.NewMockIRepository(mockCtrl)
	s.mockEvent = mockEvent.NewMockEventPublisher(mockCtrl)
	s.mockEmail = mockEmail.NewMockEmailService(mockCtrl)

	s.svc = loan.NewLoanService(s.mockRepo, s.mockEvent, s.mockEmail)
}

func TestRunServiceBillerSuite(t *testing.T) {
	suite.Run(t, new(LoanSuite))
}

func (s *LoanSuite) TestCreateLoanRequestReturnSuccess() {
	loanReq := loan.Loan{
		BorrowerID:      1,
		PrincipalAmount: 100.00,
		Rate:            5.0,
		LoanTerm:        12,
	}

	var err error
	expectedResult := err

	s.mockRepo.EXPECT().IsBorrowerExist(gomock.Any(), gomock.Any()).Return(true, nil)
	s.mockRepo.EXPECT().CreateLoan(gomock.Any(), gomock.Any()).Return(nil)
	s.mockEvent.EXPECT().Publish(gomock.Any()).Return(nil)

	err = s.svc.CreateLoan(s.ctx, &loanReq)
	s.Assert().NoError(err)
	s.Assert().Equal(expectedResult, err)

}

func (s *LoanSuite) TestApproveLoanReturnSuccess() {
	approveLoan := loan.ApprovalRequest{
		LoanId:           2,
		FieldValidatorID: 8,
	}

	loan := &loan.Loan{
		Status: loan.Proposed,
	}

	var err error
	expectedResult := err

	s.mockRepo.EXPECT().UpdateLoan(gomock.Any(), gomock.Any()).Return(nil)
	s.mockRepo.EXPECT().GetLoanByID(gomock.Any(), gomock.Any()).Return(loan, nil)
	s.mockEvent.EXPECT().Publish(gomock.Any()).Return(nil)

	err = s.svc.ApproveLoan(s.ctx, approveLoan)
	s.Assert().NoError(err)
	s.Assert().Equal(expectedResult, err)

}

func (s *LoanSuite) TestInvestLoanInvestedShouldReturnSuccess() {

	investRequest := loan.Investment{
		LoanID:     1,
		InvestorID: 100,
		Amount:     500000.00,
	}

	loan := &loan.Loan{
		Status:          loan.Approved,
		PrincipalAmount: 500000.00,
	}

	investorIds := []int{100}

	var err error
	expectedResult := err

	s.mockRepo.EXPECT().GetLoanByID(gomock.Any(), gomock.Any()).Return(loan, nil)
	s.mockRepo.EXPECT().IsInvestorExist(gomock.Any(), gomock.Any()).Return(true, nil)
	s.mockRepo.EXPECT().GetTotalInvestment(gomock.Any(), gomock.Any()).Return(0.00, nil)
	s.mockRepo.EXPECT().CreateInvestment(gomock.Any(), gomock.Any()).Return(nil)
	s.mockRepo.EXPECT().UpdateLoan(gomock.Any(), gomock.Any()).Return(nil)
	s.mockRepo.EXPECT().GetInvestorByLoanId(gomock.Any(), gomock.Any()).Return(investorIds, nil)
	s.mockEmail.EXPECT().SendAgreementLetter(gomock.Any()).Return(nil)

	s.mockEvent.EXPECT().Publish(gomock.Any()).Return(nil)

	err = s.svc.InvestLoan(s.ctx, investRequest)
	s.Assert().NoError(err)
	s.Assert().Equal(expectedResult, err)

}

func (s *LoanSuite) TestInvestLoanBelowPrincipalAmountShouldReturnSuccess() {

	investRequest := loan.Investment{
		LoanID:     1,
		InvestorID: 100,
		Amount:     5000.00,
	}

	loan := &loan.Loan{
		Status:          loan.Approved,
		PrincipalAmount: 500000.00,
	}

	var err error
	expectedResult := err

	s.mockRepo.EXPECT().GetLoanByID(gomock.Any(), gomock.Any()).Return(loan, nil)
	s.mockRepo.EXPECT().IsInvestorExist(gomock.Any(), gomock.Any()).Return(true, nil)
	s.mockRepo.EXPECT().GetTotalInvestment(gomock.Any(), gomock.Any()).Return(0.00, nil)
	s.mockRepo.EXPECT().CreateInvestment(gomock.Any(), gomock.Any()).Return(nil)

	s.mockEvent.EXPECT().Publish(gomock.Any()).Return(nil)

	err = s.svc.InvestLoan(s.ctx, investRequest)
	s.Assert().NoError(err)
	s.Assert().Equal(expectedResult, err)

}

func (s *LoanSuite) TestDisburseLoanShouldReturnSuccess() {
	disburse := loan.Disbursement{
		LoanID:         12,
		FieldOfficerID: 27,
	}

	loan := &loan.Loan{
		ID:              2,
		Rate:            5.00,
		TotalInterest:   25000.00,
		LoanTerm:        12,
		Status:          loan.Invested,
		PrincipalAmount: 500000.00,
	}

	var err error
	expectedResult := err

	s.mockRepo.EXPECT().GetLoanByID(gomock.Any(), gomock.Any()).Return(loan, nil)
	s.mockRepo.EXPECT().GetTotalInvestment(gomock.Any(), gomock.Any()).Return(500000.00, nil)
	s.mockRepo.EXPECT().IsBorrowerExist(gomock.Any(), gomock.Any()).Return(true, nil)
	s.mockRepo.EXPECT().UpdateLoan(gomock.Any(), gomock.Any()).Return(nil)

	s.mockEvent.EXPECT().Publish(gomock.Any()).Return(nil)

	err = s.svc.DisburseLoan(s.ctx, disburse)
	s.Assert().NoError(err)
	s.Assert().Equal(expectedResult, err)

}

func (s *LoanSuite) TestGetLoanShouldReturnSuccess() {

	var loanId int64
	loanId = 6

	DTOLoan := &loan.Loan{
		ID:              2,
		Rate:            5.00,
		TotalInterest:   25000.00,
		LoanTerm:        12,
		Status:          loan.Invested,
		PrincipalAmount: 500000.00,
	}

	detailInvestment := []loan.InvestmentDetail{
		{
			InvestorId:  2,
			TotalROI:    20000.00,
			TotalAmount: 5000000.00,
		},
	}

	expectedResult := loan.LoanDetail{
		ID:               DTOLoan.ID,
		Rate:             DTOLoan.Rate,
		PrincipalAmount:  DTOLoan.PrincipalAmount,
		TotalInterest:    DTOLoan.TotalInterest,
		LoanTerm:         DTOLoan.LoanTerm,
		Status:           loan.Invested,
		InvestmentDetail: detailInvestment,
	}

	s.mockRepo.EXPECT().GetLoanByID(gomock.Any(), gomock.Any()).Return(DTOLoan, nil)
	s.mockRepo.EXPECT().GetInvestmentDetail(gomock.Any(), gomock.Any()).Return(detailInvestment, nil)

	result, err := s.svc.GetLoanByID(s.ctx, loanId)
	s.Assert().NoError(err)
	s.Assert().Equal(expectedResult, result)

}
