package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/ksred/bank/accounts"
	"github.com/ksred/bank/payments"
	"net"
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

		// @TODO Send to server over TCP or similar
		sendToServer(text)
	}
}

func processCommand(text string) (result string) {
	// @TODO Receive this comm over TCP on server side
	// Commands are received split by tilde (~)
	// command~DATA

	cleanText := strings.Replace(text, "\n", "", -1)
	fmt.Printf("### %s ####\n", cleanText)
	command := strings.Split(cleanText, "~")

	// Check if we received a command
	if len(command) == 0 {
		fmt.Println("No command received")
		return
	}

	// Remove null termination from data
	command[len(command)-1] = string(bytes.Trim([]byte(command[len(command)-1]), "\x00"))

	switch command[0] {
	case "pain":
		// Check "help"
		if command[1] == "help" {
			fmt.Println("Format of PAIN transaction:\npain\npainType~senderAccountNumber@SenderBankNumber\nreceiverAccountNumber@ReceiverBankNumber\ntransactionAmount\n\nBank numbers may be left void if bank is local")
			return
		}
		result = payments.ProcessPAIN(command)
	case "camt":
	case "acmt":
		// Check "help"
		if command[1] == "help" {
			fmt.Println("") // @TODO Help section
			return
		}
		result = accounts.ProcessAccount(command)
	case "remt":
	case "reda":
	case "pacs":
	case "auth":
		break
	default:
		fmt.Println("No valid command received")
		break
	}

	return
}

func sendToServer(text string) {
	// Connect to this socket
	conn, _ := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	// Send to socket
	fmt.Fprintf(conn, text+"\n")
	// Listen for reply
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message from server: " + message)
}
