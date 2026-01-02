package repo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewManager(t *testing.T) {
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	if manager == nil {
		t.Error("Expected non-nil manager")
	}

	if manager.localRepo == nil {
		t.Error("Expected non-nil localRepo")
	}
}

func TestGetCommandSet_LocalFirst(t *testing.T) {
	// Create temporary directories
	tmpHome := t.TempDir()
	tmpBundled := t.TempDir()

	// Set up environment
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		}
	}()

	// Create local command set
	localPath := filepath.Join(tmpHome, ".shelldock")
	os.MkdirAll(localPath, 0755)
	localRepo := NewRepository(localPath)
	localCmdSet := &CommandSet{
		Name:        "test",
		Description: "Local version",
		Version:     "v1",
		Commands:    []Command{},
	}
	err := localRepo.SaveCommandSet(localCmdSet, "")
	if err != nil {
		t.Fatalf("Failed to save local command set: %v", err)
	}

	// Create bundled command set
	bundledRepo := NewRepository(tmpBundled)
	bundledCmdSet := &CommandSet{
		Name:        "test",
		Description: "Bundled version",
		Version:     "v1",
		Commands:    []Command{},
	}
	err = bundledRepo.SaveCommandSet(bundledCmdSet, "")
	if err != nil {
		t.Fatalf("Failed to save bundled command set: %v", err)
	}

	// Create manager with custom paths
	manager := &Manager{
		localRepo:   localRepo,
		bundledRepo: bundledRepo,
	}

	// Should get local version
	cmdSet, err := manager.GetCommandSet("test", false, "")
	if err != nil {
		t.Fatalf("GetCommandSet failed: %v", err)
	}

	if cmdSet.Description != "Local version" {
		t.Errorf("Expected local version, got %s", cmdSet.Description)
	}
}

func TestManagerListCommandSets(t *testing.T) {
	tmpLocal := t.TempDir()
	tmpBundled := t.TempDir()

	localRepo := NewRepository(tmpLocal)
	bundledRepo := NewRepository(tmpBundled)

	// Add to local
	localRepo.SaveCommandSet(&CommandSet{Name: "local1", Version: "v1", Commands: []Command{}}, "")
	localRepo.SaveCommandSet(&CommandSet{Name: "local2", Version: "v1", Commands: []Command{}}, "")

	// Add to bundled
	bundledRepo.SaveCommandSet(&CommandSet{Name: "bundled1", Version: "v1", Commands: []Command{}}, "")
	bundledRepo.SaveCommandSet(&CommandSet{Name: "local1", Version: "v1", Commands: []Command{}}, "") // Duplicate

	manager := &Manager{
		localRepo:   localRepo,
		bundledRepo: bundledRepo,
	}

	sets, err := manager.ListCommandSets()
	if err != nil {
		t.Fatalf("ListCommandSets failed: %v", err)
	}

	// Should have 3 unique sets: local1, local2, bundled1
	if len(sets) != 3 {
		t.Errorf("Expected 3 unique command sets, got %d: %v", len(sets), sets)
	}
}

func TestManagerListVersions(t *testing.T) {
	tmpLocal := t.TempDir()
	tmpBundled := t.TempDir()

	localRepo := NewRepository(tmpLocal)
	bundledRepo := NewRepository(tmpBundled)

	// Create versioned command set in local
	yamlContent := `name: test
versions:
  - version: "v1"
    commands: []
  - version: "v2"
    latest: true
    commands: []
`
	filePath := filepath.Join(tmpLocal, "test.yaml")
	os.WriteFile(filePath, []byte(yamlContent), 0644)

	manager := &Manager{
		localRepo:   localRepo,
		bundledRepo: bundledRepo,
	}

	versions, err := manager.ListVersions("test", false)
	if err != nil {
		t.Fatalf("ListVersions failed: %v", err)
	}

	if len(versions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(versions))
	}
}

