package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
)

//go:embed sto-plugin-binary
var stoBinary []byte

func main() {
	tmpFile, err := os.CreateTemp("", "sto-plugin-embedded-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temp file: %v\n", err)
		os.Exit(1)
	}
	tmp := tmpFile.Name()
	defer os.Remove(tmp)

	if _, err := tmpFile.Write(stoBinary); err != nil {
		tmpFile.Close()
		fmt.Fprintf(os.Stderr, "Error writing binary: %v\n", err)
		os.Exit(1)
	}
	if err := tmpFile.Chmod(0755); err != nil {
		tmpFile.Close()
		fmt.Fprintf(os.Stderr, "Error setting permissions: %v\n", err)
		os.Exit(1)
	}
	tmpFile.Close()

	cmd := exec.Command(tmp, "--run-strategy", "single-container")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}
