package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/shelldock/shelldock/internal/config"
	"github.com/shelldock/shelldock/internal/repo"
	"github.com/spf13/cobra"
)


var (
	showLocalFlag bool
	showVersionFlag string
)

var showCmd = &cobra.Command{
	Use:   "show [command-set-name]",
	Short: "Show commands in a command set without executing",
	Long: `Show the commands in a command set without executing them.
Useful for previewing what commands will be run.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		
		manager, err := repo.NewManager()
		handleError(err)

		// Check if version is specified in name (e.g., docker@1.0.0)
		var version string
		if idx := strings.Index(name, "@"); idx > 0 {
			version = name[idx+1:]
			name = name[:idx]
		} else {
			version = showVersionFlag
		}
		
		cmdSet, err := manager.GetCommandSet(name, showLocalFlag, version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Get platform
		platform, err := config.GetPlatform()
		if err != nil {
			platform = config.DetectPlatform()
		}
		
		// Show command set details without executing
		fmt.Printf("\n📦 Command Set: %s\n", cmdSet.Name)
		fmt.Printf("📝 Description: %s\n", cmdSet.Description)
		fmt.Printf("🔢 Version: %s\n", cmdSet.Version)
		fmt.Printf("🖥️  Platform: %s\n", platform)
		fmt.Printf("📋 Commands:\n\n")

		hasUnsupportedCommands := false
		for i, cmd := range cmdSet.Commands {
			fmt.Printf("  %d. %s\n", i+1, cmd.Description)
			
			// Show platform-specific command if available
			command := getCommandForPlatformShow(cmd, platform)
			if command != "" {
				fmt.Printf("     $ %s\n", command)
			} else {
				fmt.Printf("     ⚠️  No command available for platform '%s'\n", platform)
				if cmd.Platforms != nil {
					availablePlatforms := make([]string, 0, len(cmd.Platforms))
					for p := range cmd.Platforms {
						availablePlatforms = append(availablePlatforms, p)
					}
					fmt.Printf("     Available platforms: %s\n", strings.Join(availablePlatforms, ", "))
				}
				hasUnsupportedCommands = true
			}
			
			if cmd.SkipOnError {
				fmt.Printf("     ⚠️  (skip_on_error: true)\n")
			}
			fmt.Println()
		}

		if hasUnsupportedCommands {
			fmt.Printf("⚠️  Warning: Some commands are not available for platform '%s'\n", platform)
			fmt.Printf("   Consider changing your platform with: shelldock config set <platform>\n\n")
		}

		fmt.Printf("💡 To execute these commands, run: shelldock %s\n", name)
	},
}

// getCommandForPlatformShow returns the command for the specified platform
// Returns empty string if no command is available for the platform
func getCommandForPlatformShow(cmd repo.Command, platform string) string {
	if len(cmd.Platforms) > 0 {
		// 1. Try exact match
		if platformCmd, exists := cmd.Platforms[platform]; exists {
			return platformCmd
		}

		// 2. Linux hierarchy: if it's a distro, fall back to "linux"
		// If platform is not darwin or windows, we treat it as a potential linux distro
		if platform != "darwin" && platform != "windows" {
			if linuxCmd, exists := cmd.Platforms["linux"]; exists {
				return linuxCmd
			}
		}

		// 3. Safety for Windows/Darwin: if Platforms is defined but no match was found,
		// do NOT fall back to generic cmd.Command as it's likely Unix-specific.
		if platform == "windows" {
			return ""
		}
	}

	// 4. Fallback to generic command (only if no platforms defined or as final fallback for non-Windows)
	return cmd.Command
}

func init() {
	showCmd.Flags().BoolVarP(&showLocalFlag, "local", "l", false, "Only check local repository (skip bundled repository)")
	showCmd.Flags().StringVar(&showVersionFlag, "ver", "", "Show specific version or tag (default: latest). Can also use name@version format")
	showCmd.Flags().StringVar(&showVersionFlag, "version", "", "Show specific version or tag (default: latest) - alias for --ver")
}

