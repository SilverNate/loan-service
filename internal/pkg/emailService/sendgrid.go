package emailService

import (
	"github.com/sirupsen/logrus"
)

type SendGridEmailService struct {
	APIKey string
}

func NewSendGridEmailService(apiKey string) *SendGridEmailService {
	return &SendGridEmailService{APIKey: apiKey}
}

func (s *SendGridEmailService) SendAgreementLetter(investorId int) error {
	logrus.Printf("Sending agreement letter email to investor-id:  %d\n", investorId)
	return nil
}
