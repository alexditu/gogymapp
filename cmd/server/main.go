package main

import (
	"fmt"
	"os"
)

func main() {
	BINARY_NAME := "server"

	rootCmd := NewRootCmd(BINARY_NAME)
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %s failed to start: %s\n", BINARY_NAME, err)
		os.Exit(1)
	}
}
