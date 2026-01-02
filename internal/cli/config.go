package cli

import (
	"fmt"
	"os"

	"github.com/shelldock/shelldock/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage ShellDock configuration",
	Long:  "View or set ShellDock configuration settings",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		platform, err := config.GetPlatform()
		if err != nil {
			platform = config.DetectPlatform()
		}

		fmt.Println("ShellDock Configuration:")
		fmt.Printf("  Platform setting: %s\n", cfg.Platform)
		fmt.Printf("  Active platform: %s\n", platform)
		fmt.Printf("  Config file: ~/.shelldock/.sdrc\n")
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [platform]",
	Short: "Set platform (ubuntu, debian, centos, fedora, arch, darwin, windows, or auto)",
	Long: `Set the platform for command execution.
Linux distributions: ubuntu, debian, centos, rhel, fedora, arch, opensuse, alpine, etc.
Other platforms: darwin, windows
Use "auto" to automatically detect the platform and distribution.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		platform := args[0]
		
		// Allow any platform string, but validate common ones
		switch platform {
		case "auto":
			// This is always valid
		case "linux":
			// Warn that linux is ambiguous
			fmt.Fprintf(os.Stderr, "Warning: 'linux' is ambiguous. Consider using a specific distribution like 'ubuntu', 'centos', 'fedora', etc.\n")
			fmt.Fprintf(os.Stderr, "Detected distribution: %s\n", config.DetectLinuxDistribution())
		}

		cfg := &config.Config{
			Platform: platform,
		}

		if err := config.SaveConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Platform set to: %s\n", platform)
		switch platform {
		case "auto":
			detected := config.DetectPlatform()
			fmt.Printf("Auto-detected platform: %s\n", detected)
		case "linux":
			detected := config.DetectLinuxDistribution()
			fmt.Printf("Note: Using generic 'linux'. Detected distribution: %s\n", detected)
			fmt.Printf("Consider setting platform to '%s' for better compatibility.\n", detected)
		}
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
}

