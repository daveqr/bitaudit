package auditstamp

import (
	"encoding/json"
	hash "github.com/btcsuite/btcd/chaincfg/chainhash"
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
	Message   string    `json:"Message"`
	Signer    string    `json:"Signer"`
	Timestamp time.Time `json:"Timestamp"`
	Txid      string    `json:"Txid"`
	Status    Status    `json:"Status"`
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

func (as AuditStamp) Hash() hash.Hash {
	jsonstr, err := json.Marshal(as)
	if err != nil {
		panic(err)
	}

	return hash.DoubleHashH(jsonstr)
}
