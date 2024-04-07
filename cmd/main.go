package main

import (
	"fmt"

	tcpServer "github.com/Edu0liver/TCP-Server/tcp-server"
)

func main() {
	server := tcpServer.NewServer(":8080")

	go func() {
		for {
			select {
			case msg := <-server.Msgch:
				fmt.Printf("Received message from connection (%s): %s", msg.From, string(msg.Payload))
			}
		}
	}()

	server.Start()
}
