package main

import "testing"

func TestMainNoFlags(t *testing.T) {
	//@TODO Figure out a way to test os exit signals
	/*
		if os.Getenv("BE_PARSE_FLAGS") == "1" {
			parseFlags("")
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestMainNoFlags")
		cmd.Env = append(os.Environ(), "BE_PARSE_FLAGS=1")
		err := cmd.Run()

		e, _ := err.(*exec.ExitError)
		if e.Success() {
			t.Fatalf("When started with no flags, program should exit")
		}
	*/
}
