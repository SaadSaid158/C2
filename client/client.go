package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	serverAddr := "127.0.0.1:5000"

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("[-] Failed to connect to C2")
		return
	}
	defer conn.Close()

	fmt.Println("[+] Connected to C2 Server")

	for {
		fmt.Print("client> ")
		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		if command == "" {
			continue
		}

		_, err := conn.Write([]byte(command + "\n"))
		if err != nil {
			fmt.Println("[-] Failed to send command")
			return
		}

		response, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(response)
	}
}
