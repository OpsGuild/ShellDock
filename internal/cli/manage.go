package cli

import (
	"github.com/shelldock/shelldock/internal/tui"
	"github.com/spf13/cobra"
)

var manageCmd = &cobra.Command{
	Use:   "manage",
	Short: "Manage command sets (interactive UI)",
	Long:  "Open an interactive terminal UI to view, add, edit, and delete command sets",
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.Run(); err != nil {
			handleError(err)
		}
	},
}




