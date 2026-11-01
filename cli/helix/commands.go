// Package cli provides command-line tools for the Helix observability platform
package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"helix.io/pkg/go-sdk"
)

// CLI represents the main command-line interface
type CLI struct {
	version  string
	client   *helix.Client
	config   *CLIConfig
	commands map[string]Command
}

// CLIConfig holds CLI configuration
type CLIConfig struct {
	Endpoint string
	APIKey   string
	Tenant   string
	Output   string // json, table, csv
	Verbose  bool
}

// Command represents a CLI command
type Command struct {
	Name        string
	Description string
	Usage       string
	Flags       []Flag
	Run         func(ctx context.Context, args []string) error
}

// Flag represents a command-line flag
type Flag struct {
	Name        string
	Shorthand   string
	Description string
	Default     string
	Required    bool
}

// NewCLI creates a new CLI instance
func NewCLI(version string) *CLI {
	cli := &CLI{
		version:  version,
		config:   &CLIConfig{},
		commands: make(map[string]Command),
	}
	
	cli.registerCommands()
	return cli
}

// registerCommands registers all CLI commands
func (c *CLI) registerCommands() {
	c.commands["query"] = Command{
		Name:        "query",
		Description: "Execute a query against the Helix platform",
		Usage:       "helix query [expression] [flags]",
		Flags: []Flag{
			{Name: "tenant", Shorthand: "t", Description: "Tenant ID", Default: "default"},
			{Name: "start", Shorthand: "s", Description: "Start time (RFC3339 or relative)", Default: "1h"},
			{Name: "end", Shorthand: "e", Description: "End time (RFC3339 or relative)", Default: "now"},
			{Name: "step", Shorthand: "p", Description: "Step interval", Default: "1m"},
			{Name: "output", Shorthand: "o", Description: "Output format (json, table, csv)", Default: "table"},
			{Name: "limit", Shorthand: "l", Description: "Limit number of results", Default: "100"},
		},
		Run: c.runQuery,
	}
	
	c.commands["ingest"] = Command{
		Name:        "ingest",
		Description: "Ingest events from a file or stream",
		Usage:       "helix ingest [file] [flags]",
		Flags: []Flag{
			{Name: "tenant", Shorthand: "t", Description: "Tenant ID", Default: "default"},
			{Name: "source", Shorthand: "s", Description: "Source name", Default: "cli"},
			{Name: "batch-size", Shorthand: "b", Description: "Batch size", Default: "100"},
			{Name: "rate-limit", Shorthand: "r", Description: "Rate limit (events/sec)", Default: "1000"},
		},
		Run: c.runIngest,
	}
	
	c.commands["status"] = Command{
		Name:        "status",
		Description: "Show system status and metrics",
		Usage:       "helix status [flags]",
		Flags: []Flag{
			{Name: "output", Shorthand: "o", Description: "Output format (json, table)", Default: "table"},
		},
		Run: c.runStatus,
	}
	
	c.commands["alerts"] = Command{
		Name:        "alerts",
		Description: "Manage alerts",
		Usage:       "helix alerts [subcommand] [flags]",
		Flags: []Flag{
			{Name: "status", Shorthand: "s", Description: "Filter by status (firing, resolved, acknowledged)", Default: ""},
			{Name: "severity", Shorthand: "v", Description: "Filter by severity (critical, warning, info)", Default: ""},
			{Name: "limit", Shorthand: "l", Description: "Limit number of results", Default: "50"},
		},
		Run: c.runAlerts,
	}
	
	c.commands["dashboards"] = Command{
		Name:        "dashboards",
		Description: "Manage dashboards",
		Usage:       "helix dashboards [subcommand] [flags]",
		Flags: []Flag{
			{Name: "action", Shorthand: "a", Description: "Action (list, get, create, delete)", Default: "list"},
			{Name: "id", Shorthand: "i", Description: "Dashboard ID", Default: ""},
			{Name: "title", Shorthand: "t", Description: "Dashboard title", Default: ""},
		},
		Run: c.runDashboards,
	}
	
	c.commands["tenants"] = Command{
		Name:        "tenants",
		Description: "Manage tenants",
		Usage:       "helix tenants [subcommand] [flags]",
		Flags: []Flag{
			{Name: "action", Shorthand: "a", Description: "Action (list, create, delete)", Default: "list"},
			{Name: "name", Shorthand: "n", Description: "Tenant name", Default: ""},
			{Name: "rate-limit", Shorthand: "r", Description: "Rate limit", Default: "1000"},
			{Name: "quota", Shorthand: "q", Description: "Quota", Default: "10000"},
		},
		Run: c.runTenants,
	}
	
	c.commands["config"] = Command{
		Name:        "config",
		Description: "Manage CLI configuration",
		Usage:       "helix config [subcommand] [flags]",
		Flags: []Flag{
			{Name: "endpoint", Shorthand: "e", Description: "API endpoint", Default: ""},
			{Name: "api-key", Shorthand: "k", Description: "API key", Default: ""},
			{Name: "tenant", Shorthand: "t", Description: "Default tenant", Default: ""},
		},
		Run: c.runConfig,
	}
	
	c.commands["version"] = Command{
		Name:        "version",
		Description: "Show version information",
		Usage:       "helix version",
		Run:         c.runVersion,
	}
	
	c.commands["help"] = Command{
		Name:        "help",
		Description: "Show help information",
		Usage:       "helix help [command]",
		Run:         c.runHelp,
	}
}

// Run executes the CLI
func (c *CLI) Run() error {
	args := os.Args[1:]
	if len(args) == 0 {
		return c.runHelp(context.Background(), nil)
	}
	
	// Load config
	c.loadConfig()
	
	// Parse global flags
	args = c.parseGlobalFlags(args)
	
	// Get command
	cmdName := args[0]
	args = args[1:]
	
	cmd, ok := c.commands[cmdName]
	if !ok {
		return fmt.Errorf("unknown command: %s. Use 'helix help' for available commands", cmdName)
	}
	
	// Parse command flags
	args = c.parseCommandFlags(cmd, args)
	
	// Run command
	return cmd.Run(context.Background(), args)
}

// loadConfig loads CLI configuration
func (c *CLI) loadConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	
	configPath := home + "/.helix/config.json"
	data, err := os.ReadFile(configPath)
	if err != nil {
		return
	}
	
	var config CLIConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return
	}
	
	c.config = &config
	if c.config.Endpoint != "" {
		// Initialize client
		c.client = helix.NewClient(c.config.Endpoint, c.config.APIKey)
	}
}

// saveConfig saves CLI configuration
func (c *CLI) saveConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	configDir := home + "/.helix"
	os.MkdirAll(configDir, 0755)
	
	data, err := json.MarshalIndent(c.config, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(configDir+"/config.json", data, 0644)
}

// parseGlobalFlags parses global flags
func (c *CLI) parseGlobalFlags(args []string) []string {
	var result []string
	i := 0
	for i < len(args) {
		switch args[i] {
		case "--endpoint", "-e":
			if i+1 < len(args) {
				c.config.Endpoint = args[i+1]
				i += 2
				continue
			}
		case "--api-key", "-k":
			if i+1 < len(args) {
				c.config.APIKey = args[i+1]
				i += 2
				continue
			}
		case "--tenant", "-t":
			if i+1 < len(args) {
				c.config.Tenant = args[i+1]
				i += 2
				continue
			}
		case "--verbose", "-v":
			c.config.Verbose = true
			i++
			continue
		case "--output", "-o":
			if i+1 < len(args) {
				c.config.Output = args[i+1]
				i += 2
				continue
			}
		}
		result = append(result, args[i])
		i++
	}
	return result
}

// parseCommandFlags parses command-specific flags
func (c *CLI) parseCommandFlags(cmd Command, args []string) []string {
	// Simplified flag parsing
	return args
}

// Command implementations

func (c *CLI) runQuery(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("query expression required")
	}
	
	expression := args[0]
	
	if c.client == nil {
		return fmt.Errorf("not configured. Run 'helix config --endpoint <url> --api-key <key>'")
	}
	
	start := time.Now()
	results, err := c.client.Query(ctx, &helix.QueryOptions{
		Expression: expression,
		Tenant:     c.config.Tenant,
	})
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	
	duration := time.Since(start)
	
	// Output results
	switch c.config.Output {
	case "json":
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))
	case "csv":
		fmt.Println("metric,value,timestamp")
		for _, r := range results {
			fmt.Printf("%v,%v,%v\n", r.Metric, r.Value, r.Timestamp)
		}
	default:
		fmt.Printf("Query: %s\n", expression)
		fmt.Printf("Results: %d\n", len(results))
		fmt.Printf("Duration: %v\n", duration)
		fmt.Println("---")
		for _, r := range results {
			fmt.Printf("%v = %v @ %v\n", r.Metric, r.Value, r.Timestamp)
		}
	}
	
	return nil
}

func (c *CLI) runIngest(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("file path required")
	}
	
	filePath := args[0]
	
	if c.client == nil {
		return fmt.Errorf("not configured. Run 'helix config --endpoint <url> --api-key <key>'")
	}
	
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	count := 0
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		var event helix.Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			fmt.Printf("Error parsing line: %v\n", err)
			continue
		}
		
		if event.Tenant == "" {
			event.Tenant = c.config.Tenant
		}
		
		if err := c.client.IngestEvent(ctx, &event); err != nil {
			fmt.Printf("Error ingesting event: %v\n", err)
			continue
		}
		count++
	}
	
	fmt.Printf("Ingested %d events from %s\n", count, filePath)
	return nil
}

func (c *CLI) runStatus(ctx context.Context, args []string) error {
	if c.client == nil {
		return fmt.Errorf("not configured. Run 'helix config --endpoint <url> --api-key <key>'")
	}
	
	healthy, err := c.client.HealthCheck(ctx)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	
	metrics, err := c.client.Metrics(ctx)
	if err != nil {
		return fmt.Errorf("failed to get metrics: %w", err)
	}
	
	switch c.config.Output {
	case "json":
		data, _ := json.MarshalIndent(map[string]interface{}{
			"healthy": healthy,
			"metrics": metrics,
		}, "", "  ")
		fmt.Println(string(data))
	default:
		status := "❌ Unhealthy"
		if healthy {
			status = "✅ Healthy"
		}
		fmt.Printf("Status: %s\n", status)
		fmt.Println("---")
		fmt.Printf("Events Received:  %d\n", metrics["events_received"])
		fmt.Printf("Events Processed: %d\n", metrics["events_processed"])
		fmt.Printf("Active Tenants:   %d\n", metrics["active_tenants"])
		fmt.Printf("Memory Usage:     %d bytes\n", metrics["memory_usage_bytes"])
	}
	
	return nil
}

func (c *CLI) runAlerts(ctx context.Context, args []string) error {
	if c.client == nil {
		return fmt.Errorf("not configured. Run 'helix config --endpoint <url> --api-key <key>'")
	}
	
	alerts, err := c.client.GetAlerts(ctx)
	if err != nil {
		return fmt.Errorf("failed to get alerts: %w", err)
	}
	
	fmt.Printf("Alerts: %d\n", len(alerts))
	fmt.Println("---")
	for _, alert := range alerts {
		fmt.Printf("[%s] %s - %s\n", alert.Severity, alert.Name, alert.Status)
	}
	
	return nil
}

func (c *CLI) runDashboards(ctx context.Context, args []string) error {
	if c.client == nil {
		return fmt.Errorf("not configured. Run 'helix config --endpoint <url> --api-key <key>'")
	}
	
	dashboards, err := c.client.ListDashboards(ctx)
	if err != nil {
		return fmt.Errorf("failed to get dashboards: %w", err)
	}
	
	fmt.Printf("Dashboards: %d\n", len(dashboards))
	fmt.Println("---")
	for _, d := range dashboards {
		fmt.Printf("%s - %s\n", d.Title, d.Description)
	}
	
	return nil
}

func (c *CLI) runTenants(ctx context.Context, args []string) error {
	if c.client == nil {
		return fmt.Errorf("not configured. Run 'helix config --endpoint <url> --api-key <key>'")
	}
	
	tenants, err := c.client.ListTenants(ctx)
	if err != nil {
		return fmt.Errorf("failed to list tenants: %w", err)
	}
	
	fmt.Printf("Tenants: %d\n", len(tenants))
	fmt.Println("---")
	for _, t := range tenants {
		fmt.Printf("%s (%s) - Rate Limit: %d, Quota: %d\n", t.Name, t.ID, t.RateLimit, t.Quota)
	}
	
	return nil
}

func (c *CLI) runConfig(ctx context.Context, args []string) error {
	// Update config from flags
	if c.config.Endpoint != "" {
		fmt.Printf("Endpoint: %s\n", c.config.Endpoint)
	}
	if c.config.APIKey != "" {
		fmt.Printf("API Key: %s\n", c.config.APIKey)
	}
	if c.config.Tenant != "" {
		fmt.Printf("Tenant: %s\n", c.config.Tenant)
	}
	
	if err := c.saveConfig(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	
	fmt.Println("Configuration saved")
	return nil
}

func (c *CLI) runVersion(ctx context.Context, args []string) error {
	fmt.Printf("Helix CLI v%s\n", c.version)
	fmt.Println("Build: development")
	fmt.Println("Go version: 1.21")
	fmt.Println("Platform: linux/amd64")
	return nil
}

func (c *CLI) runHelp(ctx context.Context, args []string) error {
	if len(args) > 0 {
		// Show help for specific command
		cmdName := args[0]
		cmd, ok := c.commands[cmdName]
		if !ok {
			return fmt.Errorf("unknown command: %s", cmdName)
		}
		
		fmt.Printf("Command: %s\n", cmd.Name)
		fmt.Printf("Description: %s\n", cmd.Description)
		fmt.Printf("Usage: %s\n", cmd.Usage)
		fmt.Println("Flags:")
		for _, flag := range cmd.Flags {
			fmt.Printf("  --%s, -%s: %s (default: %s)\n", flag.Name, flag.Shorthand, flag.Description, flag.Default)
		}
		return nil
	}
	
	// Show general help
	fmt.Printf("Helix CLI v%s\n", c.version)
	fmt.Println("Usage: helix [command] [flags]")
	fmt.Println()
	fmt.Println("Available commands:")
	for _, cmd := range c.commands {
		fmt.Printf("  %-12s %s\n", cmd.Name, cmd.Description)
	}
	fmt.Println()
	fmt.Println("Global flags:")
	fmt.Println("  --endpoint, -e: API endpoint")
	fmt.Println("  --api-key, -k:  API key")
	fmt.Println("  --tenant, -t:   Default tenant")
	fmt.Println("  --output, -o:   Output format (json, table, csv)")
	fmt.Println("  --verbose, -v:  Verbose output")
	fmt.Println()
	fmt.Println("Use 'helix help [command]' for more information about a command.")
	
	return nil
}

// Additional helper methods for the CLI

func (c *CLI) GetCommand(name string) (Command, bool) {
	cmd, ok := c.commands[name]
	return cmd, ok
}

func (c *CLI) ListCommands() []string {
	names := make([]string, 0, len(c.commands))
	for name := range c.commands {
		names = append(names, name)
	}
	return names
}

// Interactive mode
func (c *CLI) RunInteractive() error {
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Println("Helix Interactive Mode")
	fmt.Println("Type 'help' for available commands, 'exit' to quit")
	
	for {
		fmt.Print("helix> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		
		if input == "exit" || input == "quit" {
			fmt.Println("Goodbye!")
			return nil
		}
		
		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}
		
		cmdName := args[0]
		cmd, ok := c.commands[cmdName]
		if !ok {
			fmt.Printf("Unknown command: %s\n", cmdName)
			continue
		}
		
		if err := cmd.Run(context.Background(), args[1:]); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

// Batch mode
func (c *CLI) RunBatch(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open batch file: %w", err)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	lineNum := 0
	
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}
		
		cmdName := args[0]
		cmd, ok := c.commands[cmdName]
		if !ok {
			fmt.Printf("Line %d: unknown command: %s\n", lineNum, cmdName)
			continue
		}
		
		if err := cmd.Run(context.Background(), args[1:]); err != nil {
			fmt.Printf("Line %d: error: %v\n", lineNum, err)
		}
	}
	
	return scanner.Err()
}
