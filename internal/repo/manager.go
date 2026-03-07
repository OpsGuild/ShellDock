package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	BundledRepoDir = "repository"
	LocalRepoDir   = ".shelldock"
)

type Manager struct {
	bundledRepo *Repository
	localRepo   *Repository
}

func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	localPath := filepath.Join(homeDir, LocalRepoDir)
	localRepo := NewRepository(localPath)

	bundledPaths := []string{
		"/usr/share/shelldock/repository",
		"/usr/local/share/shelldock/repository",
		filepath.Join(filepath.Dir(os.Args[0]), "..", "share", "shelldock", "repository"),
		filepath.Join(filepath.Dir(os.Args[0]), "repository"),
		"repository",
	}

	if runtime.GOOS == "windows" {
		programData := os.Getenv("ProgramData")
		if programData == "" {
			programData = `C:\ProgramData`
		}
		bundledPaths = append([]string{
			filepath.Join(programData, "shelldock", "repository"),
			filepath.Join(os.Getenv("ProgramFiles"), "shelldock", "repository"),
		}, bundledPaths...)
	}

	userRepoPath := filepath.Join(homeDir, ".shelldock", "repository")
	bundledPaths = append(bundledPaths, userRepoPath)

	var bundledRepo *Repository
	for _, path := range bundledPaths {
		if absPath, err := filepath.Abs(path); err == nil {
			if _, err := os.Stat(absPath); err == nil {
				bundledRepo = NewRepository(absPath)
				break
			}
		}
	}

	if bundledRepo == nil {
		fmt.Println("📦 No command repository found. Running initial sync...")
		if err := os.MkdirAll(userRepoPath, 0755); err == nil {
			bundledRepo = NewRepository(userRepoPath)
			bundledRepo.needsSync = true
		} else {
			bundledRepo = NewRepository("/dev/null")
		}
	}

	return &Manager{
		bundledRepo: bundledRepo,
		localRepo:   localRepo,
	}, nil
}

func (m *Manager) GetCommandSet(name string, preferLocal bool, version string) (*CommandSet, error) {
	if m.localRepo.Exists(name) {
		return m.localRepo.GetCommandSet(name, version)
	}

	if preferLocal {
		if version != "" {
			return nil, fmt.Errorf("command set '%s' version '%s' not found in local directory", name, version)
		}
		return nil, fmt.Errorf("command set '%s' not found in local directory", name)
	}

	if m.bundledRepo.Exists(name) {
		return m.bundledRepo.GetCommandSet(name, version)
	}

	if version != "" {
		return nil, fmt.Errorf("command set '%s' version '%s' not found in local directory or repository", name, version)
	}
	return nil, fmt.Errorf("command set '%s' not found in local directory or repository", name)
}

func (m *Manager) ListVersions(name string, preferLocal bool) ([]string, error) {
	if m.localRepo.Exists(name) {
		return m.localRepo.ListVersions(name)
	}

	if preferLocal {
		return []string{}, nil
	}

	if m.bundledRepo.Exists(name) {
		return m.bundledRepo.ListVersions(name)
	}

	return []string{}, nil
}

func (m *Manager) GetLocalRepo() *Repository {
	return m.localRepo
}

func (m *Manager) GetBundledRepo() *Repository {
	return m.bundledRepo
}

func (m *Manager) ListCommandSets() ([]string, error) {
	bundledSets, _ := m.bundledRepo.ListCommandSets()
	localSets, _ := m.localRepo.ListCommandSets()

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

type GroupedCommandSets struct {
	BundledGroups []CommandGroup
	LocalGroups   []CommandGroup
	Both          map[string]bool
}

func (m *Manager) ListCommandSetsGrouped() (*GroupedCommandSets, error) {
	bundledGroups, _ := m.bundledRepo.ListCommandSetsGrouped()
	localGroups, _ := m.localRepo.ListCommandSetsGrouped()

	bundledMap := make(map[string]bool)
	for _, g := range bundledGroups {
		for _, name := range g.Commands {
			bundledMap[name] = true
		}
	}

	localMap := make(map[string]bool)
	for _, g := range localGroups {
		for _, name := range g.Commands {
			localMap[name] = true
		}
	}

	both := make(map[string]bool)
	for name := range bundledMap {
		if localMap[name] {
			both[name] = true
		}
	}

	return &GroupedCommandSets{
		BundledGroups: bundledGroups,
		LocalGroups:   localGroups,
		Both:          both,
	}, nil
}
