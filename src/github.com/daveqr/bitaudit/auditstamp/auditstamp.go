package auditstamp

import (
	"time"
)

type Status int

const (
	StatusCreated = 1
	StatusPending = 2
)

func (s Status) String() string {
	switch s {
	case StatusCreated:
		return "CREATED"
	case StatusPending:
		return "PENDING"
	default:
		return "UNKNOWN"
	}
}

type AuditStamp struct {
	Message   string
	Signer    string
	Timestamp time.Time
	Txid      string
	Status    Status
}

func (as AuditStamp) StatusString() string {
	switch as.Status {
	case StatusCreated:
		return "Verified"
	case StatusPending:
		return "Pending"
	default:
		return "Unknown"
	}
}
