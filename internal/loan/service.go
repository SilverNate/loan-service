package loan

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"

	"loan-service/internal/events"
	"loan-service/internal/pkg/emailService"
)

type LoanService struct {
	repo         IRepository
	publisher    events.EventPublisher
	emailService emailService.EmailService
}

func NewLoanService(repo IRepository, publisher events.EventPublisher, emailService emailService.EmailService) *LoanService {
	return &LoanService{repo: repo, publisher: publisher, emailService: emailService}
}

func (s *LoanService) CreateLoan(ctx context.Context, request *Loan) error {
	isBorrowerExist, err := s.repo.IsBorrowerExist(ctx, request.BorrowerID)
	if err != nil {
		return err
	}

	if !isBorrowerExist {
		return errors.New("borrower does not exist")
	}

	if request.PrincipalAmount <= 0 || request.Rate <= 0 {
		return errors.New("invalid loan parameters")
	}

	request.Status = Proposed
	request.CreatedAt = time.Now()
	request.TotalInterest = CalculateTotalInterest(request.PrincipalAmount, request.Rate)

	if err := s.repo.CreateLoan(ctx, request); err != nil {
		return err
	}

	payload, _ := json.Marshal(request)
	event := events.LoanEvent{
		Type:    events.LoanCreated,
		LoanID:  request.ID,
		Payload: payload,
	}

	go s.publisher.Publish(event)

	return nil
}

func (s *LoanService) GetLoanByID(ctx context.Context, id int64) (response LoanDetail, err error) {
	loan, err := s.repo.GetLoanByID(ctx, id)
	if err != nil {
		return response, err
	}

	investmentDetail, err := s.repo.GetInvestmentDetail(ctx, loan.ID)
	if err != nil {
		return response, err
	}

	DTOLoans := LoanDetail{
		ID:                         loan.ID,
		BorrowerID:                 loan.BorrowerID,
		PrincipalAmount:            loan.PrincipalAmount,
		Rate:                       loan.Rate,
		TotalInterest:              loan.TotalInterest,
		AgreementLetter:            loan.AgreementLetter,
		LoanTerm:                   loan.LoanTerm,
		FieldValidatorID:           loan.FieldValidatorID,
		FieldValidatorPictureProof: loan.FieldValidatorPictureProof,
		FieldOfficerId:             loan.FieldOfficerId,
		ApprovalDate:               loan.ApprovalDate.Time.UTC(),
		DisbursementDate:           loan.DisbursementDate,
		Status:                     loan.Status,
		CreatedAt:                  loan.CreatedAt,
		UpdatedAt:                  loan.UpdatedAt,
		InvestmentDetail:           investmentDetail,
	}

	return DTOLoans, nil
}

func (s *LoanService) ApproveLoan(ctx context.Context, request ApprovalRequest) error {
	loan, err := s.repo.GetLoanByID(ctx, request.LoanId)
	if err != nil {
		return err
	}

	logrus.Info("only loan status proposed can continue to be approved")

	if loan.Status == Proposed {
		loan.Status = Approved

		pictureProofValidator := fmt.Sprintf("http://gcs.com/proof-validator/loan-%d", loan.ID)
		loan.FieldValidatorPictureProof = &pictureProofValidator // we can add gcsUrl after upload to server and convert multipart to byte like in disburse loan func
		loan.FieldValidatorID = request.FieldValidatorID         // field agent we can get from database agent but for now we just set 007(james bond)
		loan.UpdatedAt = time.Now()
		loan.ApprovalDate = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

		err = s.repo.UpdateLoan(ctx, loan)
		if err != nil {
			return err
		}

		payload, errMarshal := json.Marshal(request)
		if errMarshal != nil {
			logrus.Errorf("error marshalling payload: %v", errMarshal.Error())
			return errMarshal
		}

		event := events.LoanEvent{
			Type:    events.LoanApproved,
			LoanID:  loan.ID,
			Payload: payload,
		}
		// we use goroutine for publish event cause this is not mandatory to wait
		go s.publisher.Publish(event)

	} else {
		return errors.New("loan cannot be approved from current state")
	}

	return nil
}

func (s *LoanService) InvestLoan(ctx context.Context, request Investment) error {
	loan, err := s.repo.GetLoanByID(ctx, request.LoanID)
	if err != nil {
		return err
	}

	isInventorExist, err := s.repo.IsInvestorExist(ctx, request.InvestorID)
	if err != nil {
		return err
	}

	if !isInventorExist {
		return errors.New("investor does not exist")
	}

	if loan.Status != Approved && loan.Status != Invested {
		return errors.New("loan cannot receive investments in current state")
	}

	totalInvested, err := s.repo.GetTotalInvestment(ctx, request.LoanID)
	if err != nil {
		return err
	}

	if totalInvested+request.Amount > loan.PrincipalAmount {
		return errors.New("investment exceeds loan principal amount")
	}

	request.CreatedAt = time.Now()
	rate := loan.Rate / 100
	request.ROI = CalculateROI(loan.PrincipalAmount, rate, request.Amount)

	err = s.repo.CreateInvestment(ctx, &request)
	if err != nil {
		return fmt.Errorf("failed to create investment: %v", err)
	}

	if totalInvested+request.Amount == loan.PrincipalAmount {
		loan.Status = Invested
		loan.UpdatedAt = time.Now()

		logrus.Info("hardcoded url when upload to GCS")
		urlAgreementLetter := fmt.Sprintf("https://gcs.com/agreement-letter/loan-%v", loan.ID)
		loan.AgreementLetter = &urlAgreementLetter

		if err = s.repo.UpdateLoan(ctx, loan); err != nil {
			return err
		}

		investorIds, err := s.repo.GetInvestorByLoanId(ctx, request.LoanID)
		if err != nil {
			return err
		}

		err = s.sendAgreementLettersToInvestors(investorIds)
		if err != nil {
			return err
		}
	}

	payload, err := json.Marshal(request)
	if err != nil {
		logrus.Errorf("error marshalling payload: %v", err.Error())
		return err
	}

	event := events.LoanEvent{
		Type:    events.LoanInvested,
		LoanID:  loan.ID,
		Payload: payload,
	}

	go s.publisher.Publish(event)
	return nil
}

func (s *LoanService) DisburseLoan(ctx context.Context, disbursement Disbursement) error {
	loan, err := s.repo.GetLoanByID(ctx, disbursement.LoanID)
	if err != nil {
		return err
	}
	if loan.Status != Invested {
		return errors.New("loan cannot be disbursed from current state")
	}

	logrus.Info("check if amount loan can be disbursed")
	totalInvested, err := s.repo.GetTotalInvestment(ctx, loan.ID)
	if err != nil {
		return err
	}

	if totalInvested < loan.PrincipalAmount || totalInvested > loan.PrincipalAmount {
		return errors.New("failed to disburse please validate amount")
	}

	isBorrowerExist, err := s.repo.IsBorrowerExist(ctx, loan.BorrowerID)
	if err != nil {
		return err
	}

	if !isBorrowerExist {
		return errors.New("borrower does not exist")
	}

	logrus.Info("hardcoded url when upload to GCS")
	urlAgreementLetter := fmt.Sprintf("https://gcs.com/agreement-letter/loan-%v", loan.ID)

	loan.Status = Disbursed
	loan.AgreementLetter = &urlAgreementLetter
	loan.FieldOfficerId = disbursement.FieldOfficerID
	loan.DisbursementDate = time.Now()

	err = s.repo.UpdateLoan(ctx, loan)
	if err != nil {
		return err
	}

	payload, _ := json.Marshal(disbursement)
	event := events.LoanEvent{
		Type:    events.LoanDisbursed,
		LoanID:  loan.ID,
		Payload: payload,
	}
	go s.publisher.Publish(event)

	return nil
}

func (s *LoanService) sendAgreementLettersToInvestors(investorIDs []int) error {
	if s.emailService == nil {
		logrus.Error("emailService is nil")
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(investorIDs))

	for _, investorID := range investorIDs {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if err := s.emailService.SendAgreementLetter(id); err != nil {
				errChan <- fmt.Errorf("failed to send email to investor %d: %v", id, err)
			}
		}(investorID)
	}

	wg.Wait()
	close(errChan)

	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		for _, err := range errors {
			logrus.Errorf("error send email: %v", err.Error())
		}
		return fmt.Errorf("encountered errors while sending emails")
	}

	return nil
}

func CalculateROI(principalAmount, rate, amountInvested float64) float64 {
	return (principalAmount * rate) + amountInvested
}

func CalculateTotalInterest(principalAmount, rate float64) float64 {
	return principalAmount * (rate / 100)
}
