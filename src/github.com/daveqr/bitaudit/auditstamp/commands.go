package auditstamp

import (
	"encoding/json"
	bolt "github.com/coreos/bbolt"
	blockchain "github.com/daveqr/bitaudit/blockchain"
	"log"
)

const auditBucket = "audits"
const auditDb = "bitaudit.db"

var bcCommands blockchain.Commands

func WriteMessage(stamp AuditStamp) {
	log.Println("writing to bitaudit blockchain")
	bcCommands.WriteMessageToBlockchain(stamp.Message)
}

func SaveToDb(auditStamp *AuditStamp) {

	db, err := bolt.Open(auditDb, 0600, nil)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	stampJson, _ := json.Marshal(auditStamp)

	err = db.Update(func(btx *bolt.Tx) error {
		b, err := btx.CreateBucketIfNotExists([]byte(auditBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte(auditStamp.Txid), []byte(stampJson))
		if err != nil {
			log.Panic(err)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	log.Println("Saved to audit db")
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
