package blockchain

import (
	"fmt"
	"log"
	"strconv"
)

type Commands struct{}

func (commands *Commands) Send(from, to string, amount int, message string) *Transaction {
	if !ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := NewBlockchain(from)
	defer bc.Db.Close()

	tx := NewUTXOTransaction(from, to, amount, message, bc)
	bc.MineBlock([]*Transaction{tx})

	log.Println("Sent to local block chaine Transaction ID: ", tx.ID)
	return tx
}

func (commands *Commands) WriteMessageToBlockchain(message string) *Transaction {


	tx := commands.Send("16to7DptTqSrLyz7Ujymgz587zDZHEP7dr", "17rvbeXKFJt7eKhaZWePEbeWcxkyDRqV5b", 1, message)

	log.Println("Wrote message to blockstream: " + fmt.Sprintf("Transaction ID %x:", tx.ID))

	return tx
}

func (commands *Commands) PrintChain() {
	bc := NewBlockchain("")
	defer bc.Db.Close()

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (commands *Commands) CreateWallet() string {
	wallets, _ := NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	//fmt.Printf("Your new address: %s\n", address)
	log.Println("Your new address: " + address)
	return address
}

func (commands *Commands) CreateBlockchain(address string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := CreateBlockchain(address)
	bc.Db.Close()
	fmt.Println("Done!")
}

func (commands *Commands) GetBalance(address string) int {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := NewBlockchain(address)
	defer bc.Db.Close()

	balance := 0
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := bc.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
	return balance
}

func (commands *Commands) ListAddresses() {
	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()

	log.Println("Listing addresses in local blockchain")
	for _, address := range addresses {
		log.Println(address)
	}
}
