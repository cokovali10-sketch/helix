// Package main is the entry point for the Helix CLI
package main

import (
	"fmt"
	"os"

	"helix.io/cli/helix"
)

func main() {
	cli := cli.NewCLI("0.1.0")
	
	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
