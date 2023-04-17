package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Open the file for reading
	file, err := os.Open("./client/request.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Read the contents of the file into a single line
	scanner := bufio.NewScanner(file)
	var request string
	for scanner.Scan() {
		request += scanner.Text()
	}

	fmt.Print("Please enter the port number you want to listen on: ")
	var port string
	fmt.Scanln(&port)
	// Connect to the target IP address and port
	conn, err := net.Dial("tcp", "192.168.1.125:"+port)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Send the request as a single line
	fmt.Fprintf(conn, "%s\n", request)
	fmt.Println("Request sent successfully!")

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

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

		conn.Close()

		if "end" == string(buffer[:n]) {
			break
		} else {
			fmt.Println(string(buffer[:n]))
		}
	}
}
