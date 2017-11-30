package main

import (
	"github.com/daveqr/bitaudit/blockchain"
	"log"
	"github.com/daveqr/bitaudit/server"
)

var command string
var commands blockchain.Commands

// predefined local address
var genesisAddress = "16to7DptTqSrLyz7Ujymgz587zDZHEP7dr"
var toAddress = "16to7DptTqSrLyz7Ujymgz587zDZHEP7dr"


// address 1 --
// address 2 -- moyRSL4QLzLiSyH9yb6qYfBgeqzThn2wVz
// make start
// make generate BLOCKS=100
// make getinfo
// make address2
// make sendfrom1 ADDRESS=<address2> AMOUNT=10
// make generate

func main() {
	StartBLockchain()
	server.StartServer()
}

func StartBLockchain() {
	log.Println("Starting local blockchain")

	commands.CreateWallet()
	commands.CreateBlockchain(genesisAddress)
}