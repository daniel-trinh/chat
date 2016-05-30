package main

import (
	"net"
	. "net"
	"fmt"
	"bufio"
)

type Message struct {
	Data   []byte
	Sender *Conn
}

func main() {
	fmt.Println("Launching server...")
	ln, _ := net.Listen("tcp", ":8081")

	activeClients := []*Conn{}

	addClientEvent := make(chan *Conn)
	removeClientEvent := make(chan *Conn)
	messages := make(chan Message)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			go func() {
				for {
					message, err := bufio.NewReader(conn).ReadBytes('\n')
					if err != nil || len(message) == 0 {
						removeClientEvent <- &conn
						return
					} else {
						message = append(message)
						messages <- Message{message, &conn}
					}
				}
			}()
			addClientEvent <- &conn
		}
	}()

	for {
		select {
		case newClient := <-addClientEvent:
			activeClients = append(activeClients, newClient)
			fmt.Println("Client has connected. There are ", len(activeClients), "connected clients.")
		case message := <-messages:
			fmt.Println("Message received:", string(message.Data))
			fmt.Println("Sending to ", len(activeClients) - 1, "clients")
			for _, client := range activeClients {
				if client != message.Sender {
					(*client).Write(message.Data)
				}
			}
		case clientToRemove := <-removeClientEvent:
			for i, client := range activeClients {
				if client == clientToRemove {
					activeClients = append(activeClients[:i], activeClients[i + 1:]...)
				}
			}
			fmt.Println("A client has disconnected. There are ", len(activeClients), "connected clients.")
		}
	}
}