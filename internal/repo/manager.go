package repo

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	BundledRepoDir = "repository" // Relative to executable or /usr/share/shelldock
	LocalRepoDir   = ".shelldock"
)

// Manager handles both bundled (installed) and local repositories
type Manager struct {
	bundledRepo *Repository
	localRepo   *Repository
}

// NewManager creates a new repository manager
func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	localPath := filepath.Join(homeDir, LocalRepoDir)
	localRepo := NewRepository(localPath)

	// Try to find bundled repository (installed with package)
	// Check common installation paths
	bundledPaths := []string{
		"/usr/share/shelldock/repository",  // Standard Linux location
		"/usr/local/share/shelldock/repository", // Alternative location
		filepath.Join(filepath.Dir(os.Args[0]), "..", "share", "shelldock", "repository"), // Relative to executable
		filepath.Join(filepath.Dir(os.Args[0]), "repository"), // Same directory as executable
		"repository", // Current directory (for development)
	}

	var bundledRepo *Repository
	for _, path := range bundledPaths {
		if absPath, err := filepath.Abs(path); err == nil {
			if _, err := os.Stat(absPath); err == nil {
				bundledRepo = NewRepository(absPath)
				break
			}
		}
	}

	// If no bundled repo found, create a dummy one (won't find anything but won't crash)
	if bundledRepo == nil {
		bundledRepo = NewRepository("/dev/null") // Dummy path that won't exist
	}

	return &Manager{
		bundledRepo: bundledRepo,
		localRepo:   localRepo,
	}, nil
}

// GetCommandSet retrieves a command set, checking local first, then bundled repository
// version can be empty (latest), "latest", or a specific version/tag
// preferLocal flag is kept for backward compatibility but doesn't change behavior
func (m *Manager) GetCommandSet(name string, preferLocal bool, version string) (*CommandSet, error) {
	// Always check local first, then bundled (regardless of preferLocal flag)
	if m.localRepo.Exists(name) {
		return m.localRepo.GetCommandSet(name, version)
	}
	if m.bundledRepo.Exists(name) {
		return m.bundledRepo.GetCommandSet(name, version)
	}

	if version != "" {
		return nil, fmt.Errorf("command set '%s' version '%s' not found in local directory or repository", name, version)
	}
	return nil, fmt.Errorf("command set '%s' not found in local directory or repository", name)
}

// ListVersions returns all available versions for a command set
// Checks local first, then bundled repository
// preferLocal: if true, only check local (skip bundled)
func (m *Manager) ListVersions(name string, preferLocal bool) ([]string, error) {
	// Check local first
	if m.localRepo.Exists(name) {
		return m.localRepo.ListVersions(name)
	}
	
	// If --local flag is set, don't check bundled
	if preferLocal {
		return []string{}, nil
	}
	
	// Check bundled repository
	if m.bundledRepo.Exists(name) {
		return m.bundledRepo.ListVersions(name)
	}
	
	return []string{}, nil
}

// GetLocalRepo returns the local repository
func (m *Manager) GetLocalRepo() *Repository {
	return m.localRepo
}

// GetBundledRepo returns the bundled repository
func (m *Manager) GetBundledRepo() *Repository {
	return m.bundledRepo
}

// ListCommandSets returns command sets from both repositories
func (m *Manager) ListCommandSets() ([]string, error) {
	bundledSets, _ := m.bundledRepo.ListCommandSets()
	localSets, _ := m.localRepo.ListCommandSets()

	// Combine and deduplicate
	allSets := make(map[string]bool)
	for _, name := range bundledSets {
		allSets[name] = true
	}
	for _, name := range localSets {
		allSets[name] = true
	}

	var result []string
	for name := range allSets {
		result = append(result, name)
	}

	return result, nil
}

