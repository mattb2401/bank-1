package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	parseFlags()
}

func parseFlags() {
	modeFlag := flag.String("mode", "", "Test to run")

	flag.Parse()

	// Dereference
	flagParsed := *modeFlag

	switch flagParsed {
	case "client":
		// Run client for bank system
		runClient("tls")
		os.Exit(0)
		break
	case "clientNoTLS":
		// Run client for bank system
		runClient("no-tls")
		os.Exit(0)
		break
	case "server":
		// Run server for bank system
		for {
			runServer("tls")
		}
		break
	case "serverNoTLS":
		// Run server for bank system
		for {
			runServer("no-tls")
		}
		break
	default:
		fmt.Println("No valid option chosen")
		os.Exit(1)
		break
	}
}
