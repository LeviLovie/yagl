package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

	fmt.Fprintf(conn, "%s\n", what)
	conn.Close()
}

func getIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	for _, addr := range addrs {
		// Check if the address is not a loopback address and is IPv4 or IPv6
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			} else if ipnet.IP.To16() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

func input(input, sender string) {
	var req request
	err := json.Unmarshal([]byte(input), &req)
	if err != nil {
		fmt.Println("\033[31mCan't unmarshal request:", err.Error())
		fmt.Println("Request: " + input + "\033[0m")
		send(sender, "\033[31mCan't understand request\033[0m")
		return
	}

	if req.Type == "run" {
		fmt.Println("Running command: " + req.Data)
		send(sender, "\033[32mrun\033[0m:\033[34mok\033[0m:")

		command := strings.Split(req.Data, " ")
		cmd := exec.Command(command[0], command[1:]...)

		output, err := cmd.Output()
		if err != nil {
			send(sender, "\033[32mrun\033[0m:\033[34merr\033[0m:"+err.Error()+"\033[0m")
			fmt.Println("\033[31mError:", err.Error()+"\033[0m")
		}

		result := strings.Split(string(output), "\n")
		for i, line := range result {
			if line != "" {
				send(sender, "\033[32mrun\033[0m:\033[34m"+strconv.Itoa(i)+"\033[0m:\033[33m"+line+"\033[0m")
			}
		}
		send(sender, "\033[32mrun\033[0m:\033[34mdone\033[0m:")
	}

	if req.Type == "get" {
		if req.Data == "time" {
			fmt.Println("Sending time")
			now := time.Now()
			isoTime := now.Format(time.RFC3339)
			send(sender, "\033[32mget\033[0m:\033[34mtime\033[0m:\033[33m"+isoTime+"\033[0m")
			send(sender, "\033[32mget\033[0m:\033[34mdone\033[0m:")
		} else if req.Data == "ip" {
			fmt.Println("Sending IP")
			ip := getIP()
			if ip == "" {
				send(sender, "\033[32mget\033[0m:\033[34mip\033[0m:\033[31mCan't get IP\033[0m")
				fmt.Println("\033[31mCan't get IP\033[0m")
			}
			send(sender, "\033[32mget\033[0m:\033[34mip\033[0m:\033[33m"+getIP()+"\033[0m")
			send(sender, "\033[32mget\033[0m:\033[34mdone\033[0m:")
		} else {
			fmt.Println("\033[31mUnknown get request: " + req.Data + "\033[0m")
			send(sender, "\033[32mget\033[0m:\033[34merr\033[0m:\033[31mUnknown get request: "+req.Data+"\033[0m")
			send(sender, "\033[32mget\033[0m:\033[34mdone\033[0m:")
		}
	}

	if req.Type == "fget" {
		fmt.Println("Sending file: " + req.Data)
		file, err := os.Open(req.Data)
		if err != nil {
			send(sender, "\033[32mfget\033[0m:\033[34merr\033[0m:\033[31m"+err.Error()+"\033[0m")
			fmt.Println("\033[31mCan't open file:", err.Error()+"\033[0m")
		}
		defer file.Close()

		send(sender, "\033[32mfget\033[0m:\033[34mok\033[0m:")
		scanner := bufio.NewScanner(file)
		var i int
		for scanner.Scan() {
			send(sender, "\033[32mfget\033[0m:\033[34m"+strconv.Itoa(i)+"\033[0m:\033[33m"+scanner.Text()+"\033[0m")
			i++
		}
		send(sender, "\033[32mfget\033[0m:\033[34mdone\033[0m:")
	}

	send(sender, "done")
}

func main() {
	args := os.Args[1:]
	var port string
	if len(args) == 0 {
		fmt.Println("Please specify a port number")
		return
	} else {
		fmt.Println("Listening on port: \033[34m" + args[0] + "\033[0m")
		port = args[0]
	}

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
