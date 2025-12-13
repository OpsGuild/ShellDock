package cli

import (
	"fmt"

	"github.com/shelldock/shelldock/internal/repo"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync command sets from cloud repository",
	Long:  "Download and update command sets from the cloud repository",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔄 Syncing from cloud repository...")
		
		_, err := repo.NewManager()
		handleError(err)

		// TODO: Implement actual HTTP sync when cloud repo is available
		// For now, this is a placeholder
		fmt.Println("✅ Sync complete (cloud repository not yet implemented)")
		fmt.Println("💡 You can add commands locally using 'shelldock manage'")
	},
}

