package main

import (
	"time"

	"github.com/gsarmaonline/goiter/testsuite"
)

func main() {
	go testsuite.StartServer()  // Start the server in a new goroutine
	time.Sleep(2 * time.Second) // Wait for the server to start
	testsuite.Run()
}
