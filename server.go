package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":9000")

	if err != nil {
		fmt.Println("Error setting up listener", err)
		os.Exit(1)
	}

	defer listener.Close()

	fmt.Println("Server is listening on port 9000")

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting connection", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Connected to", conn.RemoteAddr().String())

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)

		if err != nil {
			fmt.Println("Error reading data", err)
			return
		}

		fmt.Println("Received message:", string(buffer[:n]))
	}
}

