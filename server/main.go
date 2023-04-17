package main

import (
	"fmt"
	"net"
)

func input(input string) {

}

func main() {
	fmt.Print("Please enter the port number you want to listen on: ")
	var port string
	fmt.Scanln(&port)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	var inputString string
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}

		fmt.Print("Request accepted from: Esc[1m " + conn.RemoteAddr().String() + "Esc[0m")

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err.Error())
			continue
		}

		inputString = string(buffer[:n])
		input(inputString)

		conn.Close()
	}
}
