package emailService

//go:generate mockgen -destination=mocks/mock.go -package=mocks -source=interface.go
type EmailService interface {
	SendAgreementLetter(investorId int) error
}
