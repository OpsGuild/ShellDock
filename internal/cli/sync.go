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

// autoSyncIfNeeded checks if the bundled repo needs initial sync and runs it.
func autoSyncIfNeeded(manager *repo.Manager) {
	bundledRepo := manager.GetBundledRepo()
	if bundledRepo == nil || !bundledRepo.NeedsSync() {
		return
	}

	bundledPath := bundledRepo.GetPath()
	if bundledPath == "" || bundledPath == "/dev/null" {
		return
	}

	count, err := syncRepository(bundledPath)
	if err != nil {
		fmt.Printf("⚠️  Auto-sync failed: %v\n", err)
		fmt.Println("💡 Run 'shelldock sync' manually to download command sets")
		return
	}

	fmt.Printf("✅ Initial sync complete! Downloaded %d command set(s)\n\n", count)
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

		bundledPath := bundledRepo.GetPath()
		if bundledPath == "" || bundledPath == "/dev/null" {
			fmt.Println("❌ Error: Bundled repository not found")
			fmt.Println("💡 Make sure ShellDock is properly installed")
			return
		}

		tempFile := filepath.Join(bundledPath, ".sync_test")
		f, err := os.Create(tempFile)
		if err != nil {
			fmt.Printf("❌ Error: No write permission for bundled repository at %s\n", bundledPath)
			fmt.Println("💡 You must run with sudo to update the bundled repository: sudo shelldock sync")
			return
		}
		_ = f.Close()
		_ = os.Remove(tempFile)

		count, err := syncRepository(bundledPath)
		if err != nil {
			fmt.Printf("❌ Error syncing repository: %v\n", err)
			return
		}

		fmt.Printf("✅ Sync complete! Updated %d command set(s) in %s\n", count, bundledPath)
	},
}

// syncRepository downloads all .yaml files from the GitHub repository.
func syncRepository(repoPath string) (int, error) {
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		return 0, fmt.Errorf("failed to create repository directory: %w", err)
	}

	count, err := processDirectory("repository", repoPath)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// processDirectory recursively processes a GitHub directory and downloads all .yaml files.
func processDirectory(dirPath, localBasePath string) (int, error) {
	url := fmt.Sprintf("%s/contents/%s", githubAPI, dirPath)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch directory listing: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return 0, fmt.Errorf("failed to fetch directory: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	var contents []githubContent
	if err := json.Unmarshal(body, &contents); err != nil {
		return 0, fmt.Errorf("failed to parse JSON: %w", err)
	}

	count := 0
	for _, item := range contents {
		if strings.Contains(item.Path, "test.yaml") {
			continue
		}

		if item.Type == "file" && strings.HasSuffix(item.Path, ".yaml") {
			relPath := strings.TrimPrefix(item.Path, "repository/")
			localPath := filepath.Join(localBasePath, relPath)

			if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
				return count, fmt.Errorf("failed to create directory: %w", err)
			}

			fileURL := fmt.Sprintf("%s/%s", githubRaw, item.Path)
			if err := downloadFile(fileURL, localPath); err != nil {
				fmt.Printf("⚠️  Warning: Could not download %s: %v\n", relPath, err)
				continue
			}

			fmt.Printf("  📥 Downloaded %s\n", relPath)
			count++
		} else if item.Type == "dir" {
			subCount, err := processDirectory(item.Path, localBasePath)
			if err != nil {
				return count, err
			}
			count += subCount
		}
	}

	return count, nil
}

// downloadFile downloads a file from a URL to a local path.
func downloadFile(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	file, err := os.Create(filePath)
	if err != nil {
		_ = resp.Body.Close()
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = io.Copy(file, resp.Body)
	_ = resp.Body.Close()
	return err
}
