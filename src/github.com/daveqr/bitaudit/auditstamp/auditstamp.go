package auditstamp

import (
	"encoding/json"
	hash "github.com/btcsuite/btcd/chaincfg/chainhash"
	"time"
)

type Status int

const (
	StatusCreated  = 1
	StatusPending  = 2
	StatusVerified = 3
)

func (s Status) String() string {
	switch s {
	case StatusCreated:
		return "CREATED"
	case StatusPending:
		return "PENDING"
	case StatusVerified:
		return "VERIFIED"
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

func (as *AuditStamp) StatusString() string {
	switch as.Status {
	case StatusCreated:
		return "Created"
	case StatusPending:
		return "Pending"
	case StatusVerified:
		return "Verified"
	default:
		return "Unknown"
	}
}

func (as *AuditStamp) Hash() hash.Hash {
	jsonstr, err := json.Marshal(as)
	if err != nil {
		panic(err)
	}

	return hash.DoubleHashH(jsonstr)
}

func (as AuditStamp) Key() string {
	return as.Signer + as.Txid
}

func (as *AuditStamp) Json() []byte {
	var tmp struct {
		Signer    string    `json:"Signer"`
		Timestamp time.Time `json:"Timestamp"`
		Txid      string    `json:"Txid"`
		Status    Status    `json:"Status"`
		Hash      hash.Hash `json:"Hash"`
	}

	tmp.Signer = as.Signer
	tmp.Timestamp = as.Timestamp
	tmp.Txid = as.Txid
	tmp.Status = as.Status
	tmp.Hash = as.Hash()
	asJson, _ := json.Marshal(tmp)

	return asJson
}

func (as *AuditStamp) ReturnJson() ([]byte, error) {
	var tmp struct {
		Txid   string `json:"Txid"`
		Status string `json:"Status"`
	}
	tmp.Txid = as.Txid
	tmp.Status = as.StatusString()

	a, err := json.Marshal(tmp)

	return a, err
}
