package auditstamp

import (
	"encoding/json"
	bolt "github.com/coreos/bbolt"
	"github.com/daveqr/bitaudit/blockchain"
	"log"
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

	err = db.Update(func(btx *bolt.Tx) error {
		b, err := btx.CreateBucketIfNotExists([]byte(auditBucket))
		if err != nil {
			log.Println(err)
			log.Panic(err)
		}

		err = b.Put([]byte(as.Key()), as.Json())
		if err != nil {
			log.Println(err)
			log.Panic(err)
		}

		return nil
	})

	if err != nil {
		log.Println(err)
		log.Panic(err)
	}

	log.Println("Saved to audit db: " + string(as.Json()))
}

func GetFromDb(key string) AuditStamp {
	db, err := bolt.Open(auditDb, 0600, nil)
	defer db.Close()
	if err != nil {
		log.Panic(err)
	}

	var sj AuditStamp
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(auditBucket))
		v := b.Get([]byte(key))

		err := json.Unmarshal(v, &sj)
		if err != nil {
			log.Panic(err)
		}

		log.Println("Returning from the db: %s", sj)

		return nil
	})

	return sj
}
