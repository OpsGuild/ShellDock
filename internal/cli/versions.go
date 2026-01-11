package cli

import (
	"fmt"
	"os"

	"github.com/shelldock/shelldock/internal/repo"
	"github.com/spf13/cobra"
)

var versionsCmd = &cobra.Command{
	Use:   "versions [command-set-name]",
	Short: "List all available versions for a command set",
	Long:  "List all available versions and tags for a specific command set",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		
		manager, err := repo.NewManager()
		handleError(err)

		versions, err := manager.ListVersions(name, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(versions) == 0 {
			fmt.Printf("No versions found for command set '%s'\n", name)
			return
		}

		fmt.Printf("Available versions for '%s':\n\n", name)
		for _, version := range versions {
			if version == "latest" {
				fmt.Printf("  * %s (default)\n", version)
			} else {
				fmt.Printf("  - %s\n", version)
			}
		}
		fmt.Printf("\nUse 'shelldock %s@<version>' or 'shelldock %s --version <version>' to run a specific version or tag\n", name, name)
	},
}

func init() {
	// This will be added in root.go
}




