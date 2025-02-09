package events

//go:generate mockgen -destination=mocks/mock.go -package=mocks -source=interface.go
type EventPublisher interface {
	Publish(event LoanEvent) error
}
