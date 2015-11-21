package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"os"
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

func sendToServer(text string) {
	// Connect to this socket
	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		log.Fatal(err)
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT, &config)

	if err != nil {
		log.Fatal(err)
	}

	// Send to socket
	fmt.Fprintf(conn, text+"\n")
	// Listen for reply
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message from server: " + message)
}
