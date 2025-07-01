package main

import "github.com/gsarmaonline/goiter/client"

func main() {
	go client.StartServer() // Start the server in a new goroutine
	client.Run()
}
