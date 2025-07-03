package main

import (
	"time"

	"github.com/gsarmaonline/goiter/client"
)

func main() {
	go client.StartServer()     // Start the server in a new goroutine
	time.Sleep(2 * time.Second) // Wait for the server to start
	client.Run()
}
