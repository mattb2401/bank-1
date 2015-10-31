package main

import (
	"bufio"
	"fmt"
	"github.com/ksred/bank/payments"
	"os"
	"strings"
)

func runClient() {

	fmt.Println("Go Banking Client\nWelcome\n")
	// We create a loop which waits for inut on std io
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		// Check if we exit
		if text == "exit\n" {
			os.Exit(0)
		}

		processCommand(text)
	}
}

func processCommand(text string) {
	// Commands are received split by tilde (~)
	// command~DATA

	cleanText := strings.Replace(text, "\n", "", -1)
	command := strings.Split(cleanText, "~")

	// Check if we received a command
	if len(command) == 0 {
		fmt.Println("No command received")
		return
	}

	switch command[0] {
	case "pain":
		// Check "help"
		if command[1] == "help" {
			fmt.Println("Format of PAIN transaction:\npain\nsenderAccountNumber@SenderBankNumber\nreceiverAccountNumber@ReceiverBankNumber\ntransactionAmount\n\nBank numbers may be left void if bank is local")
			return
		}
		payments.ProcessPAIN(command)
	case "camt":
	case "acmt":
	case "remt":
	case "reda":
	case "pacs":
	case "auth":
		break
	default:
		fmt.Println("No valid command received")
		break
	}
}
