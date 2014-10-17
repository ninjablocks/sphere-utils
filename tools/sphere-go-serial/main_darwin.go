package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("/bin/sh", "-c", "system_profiler SPHardwareDataType | sed -n 's/.*Serial Number (system).*: /OSX/p'")
	bytes, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed to wait for process because %v", err)
	}
	os.Stdout.Write(bytes[0 : len(bytes)-1])
	os.Exit(0)
}
