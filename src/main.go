package main

import (
	"./commands"
	"flag"
	"fmt"
	"strings"
)

var command string

func main() {
	var config commands.Config
	config.Password = "123"

	config.Url = "http://localhost:19001" // node 1
	config.Username = "admin1"

	//config.Url = "http://localhost:19011" // node 2
	//config.Username = "admin2"

	command := flag.String("command", "balance", "a command to run")

	var txid = "def36856aa9fbebd028f12b7aed6d1f33b26758dae141adc79699fae052d6534"
	var x = ""
	x = "gettransaction"
	//x = "getbalance"
	//x = "listtransactions"
	command = &x

	flag.Parse()

	if strings.Compare(*command, "getbalance") == 0 {
		balance, err := commands.GetBalance(config)
		if err != nil {
			panic(err)
		}
		fmt.Println(balance)
	} else if strings.Compare(*command, "gettransaction") == 0 {
		result, err := commands.GetTransaction(config, txid)
		if err != nil {
			panic(err)
		}
		fmt.Println("xxxxxxxxxxx")
		fmt.Println(result)
	} else if strings.Compare(*command, "listtransactions") == 0 {
		balance, err := commands.ListTransactions(config)
		if err != nil {
			panic(err)
		}
		fmt.Println(balance)
	/*} else if strings.Compare(*command, "signmessage") == 0 {
		balance, err := commands.SignMessage(config, txid)
		if err != nil {
			panic(err)
		}
		fmt.Println(balance)*/
	} else {
		fmt.Println("Please enter a command.")
	}
}
