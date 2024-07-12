package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")

	if err != nil {
		fmt.Println("Error connecting to server", err)
		os.Exit(1)
	}

	defer conn.Close()

	fmt.Println("Connected to server on localhost:9000")

	message := "Hello from client"

	conn.Write([]byte(message))
}
