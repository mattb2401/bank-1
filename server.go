package main

import (
	"bytes"
	"fmt"
	"github.com/ksred/bank/accounts"
	"github.com/ksred/bank/appauth"
	"github.com/ksred/bank/payments"
	"net"
	"os"
	"strings"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func runServer() {
	// @TODO Use TLS http://stackoverflow.com/questions/22666163/golang-tls-with-selfsigned-certificate
	// https://github.com/nareix/tls-example

	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	s := string(buf[:])
	// Process
	result := processCommand(s)

	// Send a response back to person contacting us.
	conn.Write([]byte(result + "\n"))
	// Close the connection when you're done with it.
	conn.Close()

}

func processCommand(text string) (result string) {
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

	// Check application auth. This is always the first value, if no token a 0 is sent
	if command[0] != "0" {
		res := appauth.CheckToken(command[0])
		if !res {
			result = "0~Incorrect token"
			return
		}
	}

	switch command[1] {
	case "appauth":
		// Check "help"
		if command[2] == "help" {
			fmt.Println("Format of appauth: appauth~userName~password")
			return
		}
		result = appauth.ProcessAppAuth(command)
		break
	case "pain":
		// Check "help"
		if command[2] == "help" {
			fmt.Println("Format of PAIN transaction:\npain\npainType~senderAccountNumber@SenderBankNumber\nreceiverAccountNumber@ReceiverBankNumber\ntransactionAmount\n\nBank numbers may be left void if bank is local")
			return
		}
		result = payments.ProcessPAIN(command)
	case "camt":
	case "acmt":
		// Check "help"
		if command[2] == "help" {
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
