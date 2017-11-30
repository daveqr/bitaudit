package auditstamp

import (
	"encoding/json"
	hash "github.com/btcsuite/btcd/chaincfg/chainhash"
	bolt "github.com/coreos/bbolt"
	blockchain "github.com/daveqr/bitaudit/blockchain"
	"log"
	"time"
)

const auditBucket = "audits"
const auditDb = "bitaudit.db"

var bcCommands blockchain.Commands

func WriteMessage(message string) {
	log.Println("writing to bitaudit blockchain")
	bcCommands.WriteMessageToBlockchain("")
}

func SaveToDb(as *AuditStamp) {
	db, err := bolt.Open(auditDb, 0600, nil)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	var tmp struct {
		Signer    string `json:"Signer"`
		Timestamp time.Time `json:"Timestamp"`
		Txid      string `json:"Txid"`
		Status    Status `json:"Status"`
		Hash      hash.Hash `json:"Hash"`
	}
	tmp.Signer = as.Signer
	tmp.Timestamp = as.Timestamp
	tmp.Txid = as.Txid
	tmp.Status = as.Status
	tmp.Hash = as.Hash()
	stampJson, _ := json.Marshal(tmp)

	err = db.Update(func(btx *bolt.Tx) error {
		b, err := btx.CreateBucketIfNotExists([]byte(auditBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte(as.Txid), stampJson)
		if err != nil {
			log.Panic(err)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	log.Println("Saved to audit db: " + string(stampJson))
}

func GetFromDb(txId string) AuditStamp {
	db, err := bolt.Open(auditDb, 0600, nil)
	defer db.Close()
	if err != nil {
		log.Panic(err)
	}

	var sj AuditStamp
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(auditBucket))
		v := b.Get([]byte(txId))

		err := json.Unmarshal(v, &sj)
		if err != nil {
			log.Panic(err)
		}

		log.Println("Returning from the db: %s", sj)

		return nil
	})

	return sj
}
