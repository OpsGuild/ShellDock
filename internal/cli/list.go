package cli

import (
	"fmt"

	"github.com/shelldock/shelldock/internal/repo"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available command sets",
	Long:  "List all available command sets from bundled repository and local directory",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := repo.NewManager()
		handleError(err)

		allSets, err := manager.ListCommandSets()
		handleError(err)

		if len(allSets) == 0 {
			fmt.Println("No command sets found.")
			return
		}

		bundledSets, _ := manager.GetBundledRepo().ListCommandSets()
		localSets, _ := manager.GetLocalRepo().ListCommandSets()

		// Create maps for quick lookup
		bundledMap := make(map[string]bool)
		for _, name := range bundledSets {
			bundledMap[name] = true
		}

		localMap := make(map[string]bool)
		for _, name := range localSets {
			localMap[name] = true
		}

		var bundledOnly []string
		var localOnly []string
		var both []string

		for _, name := range allSets {
			inBundled := bundledMap[name]
			inLocal := localMap[name]

			if inBundled && inLocal {
				both = append(both, name)
			} else if inBundled {
				bundledOnly = append(bundledOnly, name)
			} else if inLocal {
				localOnly = append(localOnly, name)
			}
		}

		if len(bundledOnly) > 0 || len(both) > 0 {
			fmt.Println("📦 Repository (bundled with installation):")
			for _, name := range bundledOnly {
				fmt.Printf("  • %s\n", name)
			}
			for _, name := range both {
				fmt.Printf("  • %s (also in local)\n", name)
			}
			fmt.Println()
		}

		if len(localOnly) > 0 {
			fmt.Println("💾 Local Repository (~/.shelldock):")
			for _, name := range localOnly {
				fmt.Printf("  • %s\n", name)
			}
		}
	},
}

