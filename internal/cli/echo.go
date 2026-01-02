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
	echoLocalFlag   bool
	echoVersionFlag string
	echoSkipFlag    string
	echoOnlyFlag    string
)

var echoCmd = &cobra.Command{
	Use:   "echo [command-set-name]",
	Short: "Echo commands in a copyable format (no descriptions or comments)",
	Long: `Echo the commands from a command set in a format that can be directly copied and pasted into a terminal.
No descriptions, comments, or formatting - just the raw commands, one per line.

Useful for:
- Copying commands to run manually
- Piping commands to a shell
- Scripting and automation

Examples:
  shelldock echo docker              # Echo all commands from docker command set
  shelldock echo docker --skip 1,2   # Echo commands skipping steps 1 and 2
  shelldock echo docker --only 3,4    # Echo only steps 3 and 4
  shelldock echo docker --ver v2     # Echo commands from version v2
  shelldock echo docker --local      # Echo from local repository only`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		manager, err := repo.NewManager()
		handleError(err)

		// Check if version is specified in name (e.g., docker@v1)
		var version string
		if idx := strings.Index(name, "@"); idx > 0 {
			version = name[idx+1:]
			name = name[:idx]
		} else {
			version = echoVersionFlag
		}

		cmdSet, err := manager.GetCommandSet(name, echoLocalFlag, version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Get platform
		platform, err := config.GetPlatform()
		if err != nil {
			platform = config.DetectPlatform()
		}

		// Filter commands if flags are provided
		commandsToRun := cmdSet.Commands
		var originalIndices []int

		if echoSkipFlag != "" || echoOnlyFlag != "" {
			var err error
			commandsToRun, originalIndices, err = filterCommands(cmdSet.Commands, echoSkipFlag, echoOnlyFlag)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			if len(commandsToRun) == 0 {
				fmt.Fprintf(os.Stderr, "Error: No commands to echo after filtering\n")
				os.Exit(1)
			}
		}

		// Echo commands in plain format (one per line, no descriptions)
		for _, cmd := range commandsToRun {
			command := getCommandForPlatform(cmd, platform)
			if command != "" {
				fmt.Println(command)
			}
			// Skip commands that don't have a command for this platform
		}
	},
}

func init() {
	echoCmd.Flags().BoolVarP(&echoLocalFlag, "local", "l", false, "Only check local repository (skip bundled repository)")
	echoCmd.Flags().StringVar(&echoVersionFlag, "ver", "", "Show specific version or tag (default: latest). Can also use name@version format")
	echoCmd.Flags().StringVar(&echoSkipFlag, "skip", "", "Skip specific steps (e.g., --skip 1,2,3 or --skip 1-3)")
	echoCmd.Flags().StringVar(&echoOnlyFlag, "only", "", "Run only specific steps (e.g., --only 1,2,3 or --only 1-3)")
}


