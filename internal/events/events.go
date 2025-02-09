package events

import (
	"encoding/json"
	"time"
)

type LoanEventType string

const (
	LoanCreated   LoanEventType = "LoanCreated"
	LoanApproved  LoanEventType = "LoanApproved"
	LoanInvested  LoanEventType = "LoanInvested"
	LoanDisbursed LoanEventType = "LoanDisbursed"
)

type LoanEvent struct {
	Type      LoanEventType   `json:"type"`
	LoanID    int64           `json:"loan_id"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}
