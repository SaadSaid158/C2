package main

import (
	"os"
	"strings"
	"time"
)

func isSandbox() bool {
	// Detect if running in a VM
	if isRunningInVM() {
		return true
	}

	// Detect if running on a single-core CPU (typical for sandboxes)
	if runtime.NumCPU() <= 1 {
		return true
	}

	// Detect if process is executing too fast (typical in sandboxes)
	start := time.Now()
	time.Sleep(2 * time.Second)
	if time.Since(start) < 2*time.Second {
		return true
	}

	return false
}

func isRunningInVM() bool {
	// Check system files for VM signatures
	data, err := os.ReadFile("/sys/class/dmi/id/product_name")
	if err == nil {
		info := strings.ToLower(string(data))
		if strings.Contains(info, "vmware") || strings.Contains(info, "virtualbox") || strings.Contains(info, "qemu") {
			return true
		}
	}

	return false
}
