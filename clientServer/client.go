package clientServer

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

var (
	// This is the FQDN from the certs generated
	CONN_HOST string
	CONN_PORT string
	CONN_TYPE string
	HTTP_PORT string
)

func initClient() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error has occurred: " + err.Error())
		os.Exit(0)
	}
	CONN_HOST = os.Getenv("CONN_HOST")
	CONN_PORT = os.Getenv("CONN_PORT")
	CONN_TYPE = os.Getenv("CONN_TYPE")
	HTTP_PORT = os.Getenv("HTTP_PORT")
}

func RunClient(mode string) {
	initClient()
	fmt.Println("Go Banking Client\nWelcome")
	// We create a loop which waits for inut on std io
	fmt.Println("Running on " + CONN_HOST + ":" + CONN_PORT)
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		// Check if we exit
		if text == "exit\n" {
			os.Exit(0)
		}

		message, err := sendToServer(text, mode)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(message)
	}
}

func sendToServer(text string, mode string) (message string, err error) {
	switch mode {
	case "tls":
		// Connect to this socket
		cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
		if err != nil {
			return "", err
		}

		config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
		conn, err := tls.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT, &config)

		if err != nil {
			return "", err
		}

		// Send to socket
		fmt.Fprintf(conn, text+"\n")

		// Listen for reply
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return "", err
		}
		return message, nil
	case "no-tls":
		// Connect to this socket
		conn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
		if err != nil {
			return "", err
		}

		// Send to socket
		fmt.Fprintf(conn, text+"\n")

		// Listen for reply
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return "", err
		}
		return message, nil
	}

	return
}
