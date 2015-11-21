package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	// This is the FQDN from the certs generated
	CONN_HOST = "bank.ksred.me"
	CONN_PORT = "6600"
	CONN_TYPE = "tcp"
)

func runServer() {
	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		log.Fatal(err)
	}

	// Load config and generate seed
	config := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.RequireAnyClientCert}
	config.Rand = rand.Reader

	// Listen for incoming connections.
	l, err := tls.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT, &config)
	if err != nil {
		log.Fatal(err)
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

	// @TODO These responses should be from the goroutine
	// Send a response back to person contacting us.
	conn.Write([]byte(result + "\n"))
	// Close the connection when you're done with it.
	conn.Close()

}
