// Package cli provides command-line tools for the Helix observability platform
package main

import (
	"fmt"
	"os"
	"strings"
)

// CLI represents the main command-line interface
type CLI struct {
	version string
	command string
	args    []string
}

// NewCLI creates a new CLI instance
func NewCLI(version string) *CLI {
	return &CLI{
		version: version,
	}
}

// Run executes the CLI
func (c *CLI) Run() error {
	if len(c.args) < 1 {
		return fmt.Errorf("no command specified. Use %s --help for available commands", c.command)
	}
	
	c.command = strings.ToLower(c.args[0])
	c.args = c.args[1:]
	
	switch c.command {
	case "query":
		return c.runQuery()
	case "ingest":
		return c.runIngest()
	case "status":
		return c.runStatus()
	case "help":
		return c.showHelp()
	default:
		return fmt.Errorf("unknown command: %s. Use %s --help for available commands", c.command, c.command)
	}
}

// runQuery executes a query against the Helix platform
func (c *CLI) runQuery() error {
	fmt.Println("Running query command")
	fmt.Printf "Helix CLI v%s - Query execution\n", c.version)
	// TODO: Implement query execution
	return nil
}

// runIngest ingests events from a file or stream
func (c *CLI) runIngest() error {
	fmt.Println("Running ingest command")
	fmt.Printf "Helix CLI v%s - Event ingestion\n", c.version)
	// TODO: Implement event ingestion
	return nil
}

// runStatus shows the status of the Helix platform
func (c *CLI) runStatus() error {
	fmt.Println("Running status command")
	fmt.Printf "Helix CLI v%s - System status\n", c.version)
	// TODO: Implement status check
	return nil
}

// showHelp displays available commands
func (c *CLI) showHelp() error {
	fmt.Printf("Helix CLI v%s\n", c.version)
	fmt.Println("Available commands:")
	fmt.Println("  query     - Execute a query against the platform")
	fmt.Println("  ingest    - Ingest events from a file or stream")
	fmt.Println("  status    - Show system status and metrics")
	fmt.Println("  help      - Show this help message")
	return nil
}

func main() {
	cli := NewCLI("0.1.0")
	
	args := os.Args[1:]
	if len(args) > 0 {
		cli.args = args
	}
	
	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}