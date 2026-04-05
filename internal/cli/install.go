package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/shelldock/shelldock/internal/repo"
	"github.com/spf13/cobra"
)

var (
	installLocalFlag   bool
	installSkipSteps   string
	installOnlySteps   string
	installVersionFlag string
	installYesFlag     bool
	installArgsFlag    string
)

var installCmd = &cobra.Command{
	Use:     "install [command-set-name]",
	Aliases: []string{"i"},
	Short:   "Install/run a command set (alias of run)",
	Long: `Install or run a saved command set.
This is an alias of "shelldock run".`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		var version string
		if idx := strings.Index(name, "@"); idx > 0 {
			version = name[idx+1:]
			name = name[:idx]
		} else {
			version = installVersionFlag
		}

		manager, err := repo.NewManager()
		handleError(err)
		autoSyncIfNeeded(manager)

		cmdSet, err := manager.GetCommandSet(name, installLocalFlag, version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		executeCommandSet(cmdSet, installSkipSteps, installOnlySteps, installYesFlag, installArgsFlag)
	},
}

func init() {
	installCmd.Flags().BoolVarP(&installLocalFlag, "local", "l", false, "Only check local repository (skip bundled repository)")
	installCmd.Flags().StringVar(&installSkipSteps, "skip", "", "Skip specific steps (comma-separated or range, e.g., 1,2,3 or 1-3)")
	installCmd.Flags().StringVar(&installOnlySteps, "only", "", "Run only specific steps (comma-separated or range, e.g., 1,3,5 or 1-3)")
	installCmd.Flags().StringVar(&installVersionFlag, "ver", "", "Run specific version or tag (uses default)")
	installCmd.Flags().StringVar(&installVersionFlag, "version", "", "Run specific version or tag (uses default) - alias for --ver")
	installCmd.Flags().BoolVarP(&installYesFlag, "yes", "a", false, "Execute all commands without prompting for confirmation")
	installCmd.Flags().StringVar(&installArgsFlag, "args", "", "Provide arguments as key=value pairs (e.g., --args name=John,email=john@example.com)")
}
