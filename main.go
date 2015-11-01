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
		runClient()
		os.Exit(0)
		break
	case "server":
		// Run server for bank system
		for {
			runServer()
		}
		break
	default:
		fmt.Println("No valid option chosen")
		os.Exit(1)
		break
	}
}
