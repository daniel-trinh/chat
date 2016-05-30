package main

import (
	"net"
	"fmt"
	"os"
	"bufio"
)

func main() {
	fmt.Println("Launching client...")
	network, err := net.Dial("tcp", ":8081")
	if err != nil {
		panic(err)
	}
	// send message
	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			message, err := reader.ReadBytes('\n')

			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			network.Write(message)
		}
	}()

	// get message
	for {
		message := make([]byte, 256)
		numBytes, err := network.Read(message)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if numBytes == 0 {
			fmt.Println("0 bytes read, connection closed")
			os.Exit(0)
		}
		fmt.Print(string(message))
	}
}