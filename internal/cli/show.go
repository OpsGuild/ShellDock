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
	showLocalFlag   bool
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
		autoSyncIfNeeded(manager)

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

		platform, err := config.GetPlatform()
		if err != nil {
			platform = config.DetectPlatform()
		}

		fmt.Printf("\n📦 Command Set: %s\n", cmdSet.Name)
		fmt.Printf("📝 Description: %s\n", cmdSet.Description)
		if cmdSet.Tag != "" {
			fmt.Printf("🔢 Version: %s @%s\n", cmdSet.Version, cmdSet.Tag)
		} else {
			fmt.Printf("🔢 Version: %s\n", cmdSet.Version)
		}
		fmt.Printf("🖥️  Platform: %s\n", platform)
		fmt.Printf("📋 Commands:\n\n")

		hasUnsupportedCommands := false
		for i, cmd := range cmdSet.Commands {
			fmt.Printf("  %d. %s\n", i+1, cmd.Description)

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

// getCommandForPlatformShow returns the command for the specified platform.
// Returns empty string if no command is available for the platform.
func getCommandForPlatformShow(cmd repo.Command, platform string) string {
	if len(cmd.Platforms) > 0 {
		if platformCmd, exists := cmd.Platforms[platform]; exists {
			return platformCmd
		}

		// Linux distros fall back to generic "linux" platform
		if platform != "darwin" && platform != "windows" {
			if linuxCmd, exists := cmd.Platforms["linux"]; exists {
				return linuxCmd
			}
		}

		// Windows with no match should not fall back to generic command (likely Unix-specific)
		if platform == "windows" {
			return ""
		}
	}

	return cmd.Command
}

func init() {
	showCmd.Flags().BoolVarP(&showLocalFlag, "local", "l", false, "Only check local repository (skip bundled repository)")
	showCmd.Flags().StringVar(&showVersionFlag, "ver", "", "Show specific version or tag (uses default). Can also use name@version format")
	showCmd.Flags().StringVar(&showVersionFlag, "version", "", "Show specific version or tag (uses default) - alias for --ver")
}
