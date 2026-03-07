package cli

import (
	"fmt"
	"sort"
	"strings"

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
		autoSyncIfNeeded(manager)

		grouped, err := manager.ListCommandSetsGrouped()
		handleError(err)

		hasBundled := len(grouped.BundledGroups) > 0
		hasLocal := len(grouped.LocalGroups) > 0

		if !hasBundled && !hasLocal {
			fmt.Println("No command sets found.")
			return
		}

		if hasBundled {
			fmt.Println("📦 Repository (bundled with installation):")
			printGroups(grouped.BundledGroups, grouped.Both, true)
			fmt.Println()
		}

		if hasLocal {
			fmt.Println("💾 Local Repository (~/.shelldock):")
			printGroups(grouped.LocalGroups, grouped.Both, false)
		}
	},
}

func printGroups(groups []repo.CommandGroup, both map[string]bool, markBoth bool) {
	sorted := make([]repo.CommandGroup, len(groups))
	copy(sorted, groups)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Name == "general" {
			return false
		}
		if sorted[j].Name == "general" {
			return true
		}
		return sorted[i].Name < sorted[j].Name
	})

	for _, group := range sorted {
		if len(group.Commands) == 0 {
			continue
		}

		cmds := make([]string, len(group.Commands))
		copy(cmds, group.Commands)
		sort.Strings(cmds)

		var parts []string
		for _, name := range cmds {
			if markBoth && both[name] {
				parts = append(parts, name+" (also in local)")
			} else {
				parts = append(parts, name)
			}
		}

		label := group.Name
		if label == "general" {
			label = "other"
		}

		fmt.Printf("  %s: %s\n", label, strings.Join(parts, ", "))
	}
}
