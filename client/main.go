package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	args := os.Args[1:]
	var port string
	var requestfile string
	var request string
	if len(args) == 0 {
		fmt.Println("Please specify a port number")
		fmt.Println("Please specify a file to send or write type and data of request")
		return
	} else if len(args) == 1 {
		port = args[0]
		fmt.Println("Please specify a file to send or write type and data of request")
		return
	} else if len(args) == 2 {
		port = args[0]
		requestfile = args[1]

		// Open the file for reading
		file, err := os.Open(requestfile)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// Read the contents of the file into a single line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			request += scanner.Text()
		}
	} else if len(args) == 3 {
		port = args[0]
		request = `{"type": "` + args[1] + `", "data": "` + args[2] + `"}`
	} else {
		fmt.Println("Too many arguments")
		return
	}

	// Connect to the target IP address and port
	conn, err := net.Dial("tcp", "192.168.1.125:"+port)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Send the request as a single line
	fmt.Fprintf(conn, "%s\n", request)
	fmt.Println("Request sent successfully!")

	fmt.Println()
	localAddr := conn.LocalAddr().String()
	_, localPort, err := net.SplitHostPort(localAddr)
	if err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", ":"+localPort)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	var exit bool
	for !exit {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err.Error())
			continue
		}

		conn.Close()

		if string(buffer[:n]) == "done\n" {
			os.Exit(0)
		}
		fmt.Print(string(buffer[:n]))
	}
}
