package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/shelldock/shelldock/internal/repo"
	"github.com/spf13/cobra"
)

var (
	rootLocalFlag bool
	rootSkipSteps string
	rootOnlySteps string
	rootVersionFlag string
	rootYesFlag bool
)

var rootCmd = &cobra.Command{
	Use:   "shelldock [command-set-name]",
	Short: "ShellDock - A repository for shell commands",
	Long: `ShellDock is a fast, cross-platform tool for managing and executing
saved shell commands from bundled repository or local directory.

You can run a command set directly:
  shelldock docker
  shelldock --local docker

Or use subcommands:
  shelldock run docker
  shelldock list
  shelldock manage`,
	Version: "1.0.0",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// If a command set name is provided, run it
		if len(args) == 1 {
			name := args[0]
			
			// Check if version is specified in name (e.g., docker@1.0.0)
			var version string
			if idx := strings.Index(name, "@"); idx > 0 {
				version = name[idx+1:]
				name = name[:idx]
			} else {
				version = rootVersionFlag
			}
			
			manager, err := repo.NewManager()
			handleError(err)

			cmdSet, err := manager.GetCommandSet(name, rootLocalFlag, version)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			executeCommandSet(cmdSet, rootSkipSteps, rootOnlySteps, rootYesFlag)
			return
		}
		// Otherwise show help
		cmd.Help()
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVarP(&rootLocalFlag, "local", "l", false, "Only check local repository (skip bundled repository)")
	rootCmd.Flags().StringVar(&rootSkipSteps, "skip", "", "Skip specific steps (comma-separated or range, e.g., 1,2,3 or 1-3)")
	rootCmd.Flags().StringVar(&rootOnlySteps, "only", "", "Run only specific steps (comma-separated or range, e.g., 1,3,5 or 1-3)")
	rootCmd.Flags().StringVar(&rootVersionFlag, "ver", "", "Run specific version or tag (default: latest). Can also use name@version format")
	rootCmd.Flags().BoolVarP(&rootYesFlag, "yes", "y", false, "Execute commands without prompting for confirmation")
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(echoCmd)
	rootCmd.AddCommand(manageCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(versionsCmd)
	rootCmd.AddCommand(configCmd)
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
