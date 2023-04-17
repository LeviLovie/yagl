package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type request struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func send(where, what string) {
	conn, err := net.Dial("tcp", where)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Send the request as a single line
	fmt.Fprintf(conn, "%s\n", what)
}

func input(input, sender string) {
	var req request
	err := json.Unmarshal([]byte(input), &req)
	if err != nil {
		fmt.Println("\033[31mCan't unmarshal request:", err.Error())
		fmt.Println("Request: " + input + "\033[0m")
		return
	}

	if req.Type == "run" {
		fmt.Println("Running command: " + req.Data)
	}

	if req.Type == "get" {
		if req.Data == "time" {
			fmt.Println("Sending time")
			now := time.Now()
			isoTime := now.Format(time.RFC3339)
			send(sender, "get:time:"+isoTime)
		}
	}

	if req.Type == "fget" {
		fmt.Println("Sending file: " + req.Data)
		file, err := os.Open(req.Data)
		if err != nil {
			send(sender, "fget:err:"+err.Error())
		}
		defer file.Close()

		send(sender, "fget:ok")
		scanner := bufio.NewScanner(file)
		var i int
		for scanner.Scan() {
			send(sender, "fget:"+strconv.Itoa(i)+":"+scanner.Text())
		}
		send(sender, "fget:done")
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
		fmt.Println()
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
		input(inputString, conn.RemoteAddr().String())

		conn.Close()
	}
}
