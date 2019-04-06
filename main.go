package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/mattb2401/bank/cmd"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error has occurred: " + err.Error())
		os.Exit(0)
	}
	cmd.Execute()
}
