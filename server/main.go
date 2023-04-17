package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type request struct {
	Sender string `json:"sender"`
	Time   string `json:"time"`
	Type   string `json:"type"`
	Data   string `json:"data"`
}

func input(input string) {
	var req request
	err := json.Unmarshal([]byte(input), &req)
	if err != nil {
		fmt.Println("\033[31mCan't unmarshal request:\033[0m", err.Error())
		fmt.Println("Request: " + input)
		return
	}

	if req.Type == "run" {
		fmt.Println("Running command: " + req.Data)
	}
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

		fmt.Println("Request accepted from: \033[34m" + conn.RemoteAddr().String() + "\033[0m")

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
