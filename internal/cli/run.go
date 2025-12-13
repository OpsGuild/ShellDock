package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/term"

	"github.com/shelldock/shelldock/internal/config"
	"github.com/shelldock/shelldock/internal/repo"
	"github.com/spf13/cobra"
)

var (
	localFlag bool
	skipSteps string
	onlySteps string
	versionFlag string
	yesFlag bool
)

// parseStepNumbers parses comma-separated step numbers (1-indexed)
func parseStepNumbers(input string) (map[int]bool, error) {
	if input == "" {
		return nil, nil
	}
	
	steps := make(map[int]bool)
	parts := strings.Split(input, ",")
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		// Check for range (e.g., "1-3")
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}
			
			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid start number in range: %s", rangeParts[0])
			}
			
			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid end number in range: %s", rangeParts[1])
			}
			
			if start > end {
				return nil, fmt.Errorf("range start (%d) must be <= end (%d)", start, end)
			}
			
			for i := start; i <= end; i++ {
				steps[i] = true
			}
		} else {
			num, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid step number: %s", part)
			}
			if num < 1 {
				return nil, fmt.Errorf("step numbers must be >= 1, got: %d", num)
			}
			steps[num] = true
		}
	}
	
	return steps, nil
}

// filterCommands filters commands based on --skip and --only flags
func filterCommands(commands []repo.Command, skipSteps, onlySteps string) ([]repo.Command, []int, error) {
	skipMap, err := parseStepNumbers(skipSteps)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid --skip format: %w", err)
	}
	
	onlyMap, err := parseStepNumbers(onlySteps)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid --only format: %w", err)
	}
	
	if skipSteps != "" && onlySteps != "" {
		return nil, nil, fmt.Errorf("cannot use both --skip and --only flags together")
	}
	
	var filtered []repo.Command
	var originalIndices []int
	
	for i, cmd := range commands {
		stepNum := i + 1 // 1-indexed
		
		// If --only is specified, only include steps in the map
		if onlySteps != "" {
			if !onlyMap[stepNum] {
				continue
			}
		}
		
		// If --skip is specified, exclude steps in the map
		if skipSteps != "" {
			if skipMap[stepNum] {
				continue
			}
		}
		
		filtered = append(filtered, cmd)
		originalIndices = append(originalIndices, stepNum)
	}
	
	return filtered, originalIndices, nil
}

// getCommandForPlatform returns the command for the specified platform
// Returns empty string if no command is available for the platform
func getCommandForPlatform(cmd repo.Command, platform string) string {
	// If platforms map exists and has entry for this platform, use it
	if cmd.Platforms != nil {
		if platformCmd, exists := cmd.Platforms[platform]; exists {
			return platformCmd
		}
		// If platforms map exists but no entry for this platform, and no fallback command
		if cmd.Command == "" {
			return "" // No command available for this platform
		}
	}

	// Fallback to generic command
	return cmd.Command
}

// executeCommandSet is the shared logic for running command sets
func executeCommandSet(cmdSet *repo.CommandSet, skipSteps, onlySteps string, yesFlag bool) {
	// Get platform
	platform, err := config.GetPlatform()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to get platform: %v, using auto-detected\n", err)
		platform = config.DetectPlatform()
	}
	
	// Filter commands if flags are provided
	commandsToRun := cmdSet.Commands
	originalIndices := make([]int, len(cmdSet.Commands))
	for i := range originalIndices {
		originalIndices[i] = i + 1
	}
	
	if skipSteps != "" || onlySteps != "" {
		var err error
		commandsToRun, originalIndices, err = filterCommands(cmdSet.Commands, skipSteps, onlySteps)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		
		if len(commandsToRun) == 0 {
			fmt.Fprintf(os.Stderr, "Error: No commands to execute after filtering\n")
			os.Exit(1)
		}
	}
	
	fmt.Printf("\n📦 Command Set: %s\n", cmdSet.Name)
	fmt.Printf("📝 Description: %s\n", cmdSet.Description)
	fmt.Printf("🔢 Version: %s\n", cmdSet.Version)
	fmt.Printf("🖥️  Platform: %s\n", platform)
	
	if skipSteps != "" {
		fmt.Printf("⏭️  Skipping steps: %s\n", skipSteps)
	} else if onlySteps != "" {
		fmt.Printf("🎯 Running only steps: %s\n", onlySteps)
	}
	
	fmt.Printf("📋 Commands to execute:\n\n")

	hasUnsupportedCommands := false
	for i, cmd := range commandsToRun {
		originalNum := originalIndices[i]
		fmt.Printf("  %d. %s\n", originalNum, cmd.Description)
		command := getCommandForPlatform(cmd, platform)
		if command == "" {
			fmt.Printf("     ⚠️  No command available for platform '%s'\n", platform)
			if cmd.Platforms != nil {
				availablePlatforms := make([]string, 0, len(cmd.Platforms))
				for p := range cmd.Platforms {
					availablePlatforms = append(availablePlatforms, p)
				}
				fmt.Printf("     Available platforms: %s\n", strings.Join(availablePlatforms, ", "))
			}
			hasUnsupportedCommands = true
		} else {
			fmt.Printf("     $ %s\n", command)
		}
		fmt.Println()
	}

	if hasUnsupportedCommands {
		fmt.Printf("⚠️  Warning: Some commands are not available for platform '%s'\n", platform)
		fmt.Printf("   Consider changing your platform with: shelldock config set <platform>\n")
		fmt.Printf("   Or use --yes flag to skip unsupported commands during execution\n\n")
	}

	// Skip prompt if --yes flag is set
	if !yesFlag {
		// Check if stdin is a terminal
		if !term.IsTerminal(int(os.Stdin.Fd())) {
			// Not a terminal (e.g., piped input), don't prompt
			fmt.Println("⚠️  Not running in a terminal. Use --yes flag to execute without prompt.")
			return
		}
		
		// Read from stdin with proper terminal handling
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Do you want to execute these commands? (y/N): ")
		
		// Ensure the prompt is displayed immediately
		os.Stdout.Sync()
		
		// Read the response - this will block until user presses Enter
		response, err := reader.ReadString('\n')
		if err != nil {
			// If we can't read (e.g., stdin is closed), cancel
			fmt.Fprintf(os.Stderr, "\nError reading input: %v\n", err)
			fmt.Println("Cancelled.")
			return
		}
		
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			fmt.Println("Cancelled.")
			return
		}
	}

	fmt.Println("\n🚀 Executing commands...\n")

	for i, cmd := range commandsToRun {
		originalNum := originalIndices[i]
		command := getCommandForPlatform(cmd, platform)
		if command == "" {
			fmt.Printf("[%d/%d] %s (step %d)\n", i+1, len(commandsToRun), cmd.Description, originalNum)
			fmt.Printf("⚠️  Skipping: No command available for platform '%s'\n\n", platform)
			continue
		}
		
		fmt.Printf("[%d/%d] %s (step %d)\n", i+1, len(commandsToRun), cmd.Description, originalNum)
		fmt.Printf("$ %s\n", command)

		execCmd := exec.Command("sh", "-c", command)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		execCmd.Stdin = os.Stdin

		if err := execCmd.Run(); err != nil {
			if cmd.SkipOnError {
				fmt.Printf("⚠️  Command failed but continuing (skip_on_error=true)\n\n")
				continue
			}
			fmt.Fprintf(os.Stderr, "\n❌ Command failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Success\n")
	}

	fmt.Println("🎉 All commands executed successfully!")
}

var runCmd = &cobra.Command{
	Use:   "run [command-set-name]",
	Short: "Run a saved command set",
	Long: `Run a saved command set. By default, searches local directory first,
       then bundled repository. Use --local or -l to only check local directory.

You can skip specific steps with --skip:
  shelldock run docker --skip 1,2,3
  shelldock run docker --skip 1-3

Or run only specific steps with --only:
  shelldock run docker --only 1,3,5
  shelldock run docker --only 1-3`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		
		// Check if version is specified in name (e.g., docker@1.0.0)
		var version string
		if idx := strings.Index(name, "@"); idx > 0 {
			version = name[idx+1:]
			name = name[:idx]
		} else {
			version = versionFlag
		}
		
		manager, err := repo.NewManager()
		handleError(err)

		cmdSet, err := manager.GetCommandSet(name, localFlag, version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		executeCommandSet(cmdSet, skipSteps, onlySteps, yesFlag)
	},
}

func init() {
		runCmd.Flags().BoolVarP(&localFlag, "local", "l", false, "Only check local repository (skip bundled repository)")
	runCmd.Flags().StringVar(&skipSteps, "skip", "", "Skip specific steps (comma-separated or range, e.g., 1,2,3 or 1-3)")
	runCmd.Flags().StringVar(&onlySteps, "only", "", "Run only specific steps (comma-separated or range, e.g., 1,3,5 or 1-3)")
	runCmd.Flags().StringVar(&versionFlag, "ver", "", "Run specific version or tag (default: latest)")
	runCmd.Flags().BoolVarP(&yesFlag, "yes", "y", false, "Execute commands without prompting for confirmation")
}

