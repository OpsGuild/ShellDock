package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/shelldock/shelldock/internal/repo"
	"github.com/spf13/cobra"
)

const (
	githubRepo = "OpsGuild/ShellDock"
	githubAPI  = "https://api.github.com/repos/" + githubRepo
	githubRaw  = "https://raw.githubusercontent.com/" + githubRepo + "/master"
)

type githubContent struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync command sets from cloud repository",
	Long:  "Download and update command sets from the cloud repository",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔄 Syncing from cloud repository...")

		manager, err := repo.NewManager()
		handleError(err)

		bundledRepo := manager.GetBundledRepo()
		if bundledRepo == nil {
			fmt.Println("❌ Error: Could not find bundled repository path")
			fmt.Println("💡 Make sure ShellDock is properly installed")
			return
		}

		// Get the actual path of the bundled repository
		bundledPath := bundledRepo.GetPath()
		if bundledPath == "" || bundledPath == "/dev/null" {
			fmt.Println("❌ Error: Bundled repository not found")
			fmt.Println("💡 Make sure ShellDock is properly installed")
			return
		}

		// Check if we have write permissions
		if _, err := os.Stat(bundledPath); err != nil {
			fmt.Printf("❌ Error: Cannot access bundled repository at %s\n", bundledPath)
			fmt.Println("💡 You may need to run with sudo to update the bundled repository")
			return
		}

		// Sync repository files
		count, err := syncRepository(bundledPath)
		if err != nil {
			fmt.Printf("❌ Error syncing repository: %v\n", err)
			return
		}

		fmt.Printf("✅ Sync complete! Updated %d command set(s)\n", count)
	},
}

// syncRepository downloads all .yaml files from the GitHub repository
func syncRepository(repoPath string) (int, error) {
	// Ensure the repository directory exists
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		return 0, fmt.Errorf("failed to create repository directory: %w", err)
	}

	// Process the repository directory recursively
	count, err := processDirectory("repository", repoPath)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// processDirectory recursively processes a GitHub directory and downloads all .yaml files
func processDirectory(dirPath, localBasePath string) (int, error) {
	url := fmt.Sprintf("%s/contents/%s", githubAPI, dirPath)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch directory listing: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to fetch directory: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	var contents []githubContent
	if err := json.Unmarshal(body, &contents); err != nil {
		return 0, fmt.Errorf("failed to parse JSON: %w", err)
	}

	count := 0
	for _, item := range contents {
		// Skip test files
		if strings.Contains(item.Path, "test.yaml") {
			continue
		}

		if item.Type == "file" && strings.HasSuffix(item.Path, ".yaml") {
			// Calculate relative path from repository/
			relPath := strings.TrimPrefix(item.Path, "repository/")
			localPath := filepath.Join(localBasePath, relPath)

			// Create subdirectory if needed
			if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
				return count, fmt.Errorf("failed to create directory: %w", err)
			}

			// Download file
			fileURL := fmt.Sprintf("%s/%s", githubRaw, item.Path)
			if err := downloadFile(fileURL, localPath); err != nil {
				fmt.Printf("⚠️  Warning: Could not download %s: %v\n", relPath, err)
				continue
			}

			fmt.Printf("  📥 Downloaded %s\n", relPath)
			count++
		} else if item.Type == "dir" {
			// Recursively process subdirectories
			subCount, err := processDirectory(item.Path, localBasePath)
			if err != nil {
				return count, err
			}
			count += subCount
		}
	}

	return count, nil
}

// downloadFile downloads a file from a URL to a local path
func downloadFile(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
