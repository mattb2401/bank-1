package main

import "testing"

func TestRunTLSServer(t *testing.T) {
	//@TODO Figure out how to exit this when successful
	/*
		errs := make(chan *bankError, 1)

		go func() {
			_, err := runServer("tls")
			if err != nil {
				errs <- err
			}
			close(errs)
			return
		}()

		for err := range errs {
			if err != nil {
				t.Errorf("RunTLSServer does not pass. Looking for %v, got %v", nil, errs)
			}
		}
	*/
}

//@TODO Implement TCP (no TLS) when exit if successful implemented

func TestHandleRequest(t *testing.T) {
	/*
		l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)

		if err != nil {
			t.Errorf("HandleRequest does not pass. Does not Listen. Looking for %v, got %v", nil, err)
		}
		conn, err := l.Accept()
		bankErr := handleRequest(conn)
		l.Close()
		if bankErr != nil {
			t.Errorf("HandleRequest does not pass. Looking for %v, got %v", nil, bankErr)
		}
	*/
}
