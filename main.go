package main

import (
	"errors"
	"flag"
	"log"
	"os"
)

const (
	// This is the FQDN from the certs generated
	CONN_HOST = "localhost"
	CONN_PORT = "3300"
	CONN_TYPE = "tcp"
)

func main() {
	modeFlag := flag.String("mode", "", "Test to run")

	flag.Parse()

	// Dereference
	flagParsed := *modeFlag

	err := parseFlags(flagParsed)
	if err != nil {
		log.Fatalf("Error starting, err: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func parseFlags(flagParsed string) (err error) {
	switch flagParsed {
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
			runServer("tls")
		}
	case "serverNoTLS":
		// Run server for bank system
		for {
			runServer("no-tls")
		}
	default:
		return errors.New("No valid option chosen. Valid options: client, clientNoTLS, server, serverNoTLS")
	}

	return
}
