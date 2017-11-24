package main

import (
	"./commands"
	"flag"
	"fmt"
	"strings"
)

var command string

func mainxx() {
	var config commands.Config
	config.Url = "http://localhost:19001/"
	config.Username = "admin1"
	config.Password = "123"

	command := flag.String("command", "balance", "a command to run")

	var txid = "52cec53b436c01f4699aaf3e8cc68988c959341cd8d9efe4c74379784e8281c6"
	var x = "gettransactions"
	command = &x

	flag.Parse()

	if strings.Compare(*command, "balance") == 0 {
		balance, err := commands.GetBalance(config)
		if err != nil {
			panic(err)
		}
		fmt.Println(balance)
	} else if strings.Compare(*command, "gettransactions") == 0 {
		result, err := commands.GetTransaction(config, txid)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	} else if strings.Compare(*command, "listtransactions") == 0 {
		balance, err := commands.ListTransactions(config)
		if err != nil {
			panic(err)
		}
		fmt.Println(balance)
	} else {
		fmt.Println("Please enter a command.")
	}
}
