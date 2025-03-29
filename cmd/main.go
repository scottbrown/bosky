// Package main implements the beacon CLI tool for sending EC2
// user data events to EventBridge
package main

import (
	"fmt"
	"os"
)

// main is the entry point for the CLI application
func main() {
	Execute()
}

// Execute run the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
