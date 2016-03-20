package main

import (
	"errors"
	"log"
	"os"
)

const (
	// This is the FQDN from the certs generated
	CONN_HOST = "localhost"
	CONN_PORT = "3300"
	CONN_TYPE = "tcp"
	HTTP_PORT = "8443"
)

func main() {
	argClientServer := os.Args[1]

	err := RunHttpServer()
	if err != nil {
		log.Fatalf("Could not start HTTP server. " + err.Error())
	}

	err = parseArguments(argClientServer)
	if err != nil {
		log.Fatalf("Error starting, err: %v\n", err)
	}
	os.Exit(0)
}

func parseArguments(arg string) (err error) {
	switch arg {
	case "client":
		// Run client for bank system
		runClient("tls")
		break
	case "clientNoTLS":
		// Run client for bank system
		runClient("no-tls")
		break
	case "server":
		// Run server for bank system
		for {
			err := RunHttpServer()
			if err != nil {
				log.Fatalf("Could not start HTTP server. " + err.Error())
			}
			runServer("tls")
		}
	case "serverNoTLS":
		// Run server for bank system
		for {
			err := RunHttpServer()
			if err != nil {
				log.Fatalf("Could not start HTTP server. " + err.Error())
			}
			runServer("no-tls")
		}
	default:
		return errors.New("No valid option chosen. Valid options: client, clientNoTLS, server, serverNoTLS")
	}

	return
}
