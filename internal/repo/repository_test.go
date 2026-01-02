package repo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewRepository(t *testing.T) {
	repo := NewRepository("/test/path")
	if repo.path != "/test/path" {
		t.Errorf("Expected path /test/path, got %s", repo.path)
	}
}

func TestExtractVersionNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"v1", 1},
		{"v2", 2},
		{"v10", 10},
		{"1", 1},
		{"V1", 1},
		{"invalid", 0},
		{"", 0},
	}

	for _, tt := range tests {
		result := extractVersionNumber(tt.input)
		if result != tt.expected {
			t.Errorf("extractVersionNumber(%q) = %d, expected %d", tt.input, result, tt.expected)
		}
	}
}

func TestGetCommandSet_SingleVersion(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewRepository(tmpDir)

	// Create a test YAML file
	yamlContent := `name: test
description: Test command set
version: "v1"
commands:
  - description: Test command
    command: echo "test"
`
	filePath := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(filePath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cmdSet, err := repo.GetCommandSet("test", "")
	if err != nil {
		t.Fatalf("GetCommandSet failed: %v", err)
	}

	if cmdSet.Name != "test" {
		t.Errorf("Expected name 'test', got %s", cmdSet.Name)
	}
	if cmdSet.Version != "v1" {
		t.Errorf("Expected version 'v1', got %s", cmdSet.Version)
	}
	if len(cmdSet.Commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(cmdSet.Commands))
	}
}

func TestGetCommandSet_Versioned(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewRepository(tmpDir)

	yamlContent := `name: test
description: Test command set
versions:
  - version: "v1"
    description: Version 1
    commands:
      - description: Command 1
        command: echo "v1"
  - version: "v2"
    latest: true
    description: Version 2
    commands:
      - description: Command 2
        command: echo "v2"
`
	filePath := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(filePath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test latest version
	cmdSet, err := repo.GetCommandSet("test", "")
	if err != nil {
		t.Fatalf("GetCommandSet failed: %v", err)
	}
	if cmdSet.Version != "v2" {
		t.Errorf("Expected latest version 'v2', got %s", cmdSet.Version)
	}

	// Test specific version
	cmdSet, err = repo.GetCommandSet("test", "v1")
	if err != nil {
		t.Fatalf("GetCommandSet failed: %v", err)
	}
	if cmdSet.Version != "v1" {
		t.Errorf("Expected version 'v1', got %s", cmdSet.Version)
	}
}

func TestGetCommandSet_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewRepository(tmpDir)

	_, err := repo.GetCommandSet("nonexistent", "")
	if err == nil {
		t.Error("Expected error for nonexistent command set")
	}
}

func TestGetCommandSet_WithArgs(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewRepository(tmpDir)

	yamlContent := `name: test
description: Test with args
version: "v1"
commands:
  - description: Test command
    command: echo "Hello {{name}}"
    args:
      - name: name
        prompt: "Enter name"
        required: true
`
	filePath := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(filePath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cmdSet, err := repo.GetCommandSet("test", "")
	if err != nil {
		t.Fatalf("GetCommandSet failed: %v", err)
	}

	if len(cmdSet.Commands[0].Args) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(cmdSet.Commands[0].Args))
	}
	if cmdSet.Commands[0].Args[0].Name != "name" {
		t.Errorf("Expected arg name 'name', got %s", cmdSet.Commands[0].Args[0].Name)
	}
}

func TestListVersions(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewRepository(tmpDir)

	yamlContent := `name: test
versions:
  - version: "v1"
    commands: []
  - version: "v2"
    latest: true
    commands: []
  - version: "v3"
    commands: []
`
	filePath := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(filePath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	versions, err := repo.ListVersions("test")
	if err != nil {
		t.Fatalf("ListVersions failed: %v", err)
	}

	if len(versions) != 3 {
		t.Errorf("Expected 3 versions, got %d", len(versions))
	}

	// Check that latest is marked
	foundLatest := false
	for _, v := range versions {
		if v == "v2 (latest)" {
			foundLatest = true
			break
		}
	}
	if !foundLatest {
		t.Error("Expected to find 'v2 (latest)' in versions")
	}
}

func TestSaveCommandSet(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewRepository(tmpDir)

	cmdSet := &CommandSet{
		Name:        "test",
		Description: "Test set",
		Version:     "v1",
		Commands: []Command{
			{
				Description: "Test command",
				Command:     "echo test",
			},
		},
	}

	err := repo.SaveCommandSet(cmdSet, "")
	if err != nil {
		t.Fatalf("SaveCommandSet failed: %v", err)
	}

	// Verify it was saved
	saved, err := repo.GetCommandSet("test", "")
	if err != nil {
		t.Fatalf("Failed to read saved command set: %v", err)
	}

	if saved.Name != "test" {
		t.Errorf("Expected name 'test', got %s", saved.Name)
	}
}

func TestListCommandSets(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewRepository(tmpDir)

	// Create multiple command sets
	cmdSets := []*CommandSet{
		{Name: "test1", Version: "v1", Commands: []Command{}},
		{Name: "test2", Version: "v1", Commands: []Command{}},
		{Name: "test3", Version: "v1", Commands: []Command{}},
	}

	for _, cmdSet := range cmdSets {
		err := repo.SaveCommandSet(cmdSet, "")
		if err != nil {
			t.Fatalf("Failed to save command set: %v", err)
		}
	}

	sets, err := repo.ListCommandSets()
	if err != nil {
		t.Fatalf("ListCommandSets failed: %v", err)
	}

	if len(sets) != 3 {
		t.Errorf("Expected 3 command sets, got %d", len(sets))
	}
}

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewRepository(tmpDir)

	if repo.Exists("test") {
		t.Error("Expected 'test' to not exist")
	}

	cmdSet := &CommandSet{Name: "test", Version: "v1", Commands: []Command{}}
	err := repo.SaveCommandSet(cmdSet, "")
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	if !repo.Exists("test") {
		t.Error("Expected 'test' to exist")
	}
}

func TestDeleteCommandSet(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewRepository(tmpDir)

	cmdSet := &CommandSet{Name: "test", Version: "v1", Commands: []Command{}}
	err := repo.SaveCommandSet(cmdSet, "")
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	err = repo.DeleteCommandSet("test")
	if err != nil {
		t.Fatalf("DeleteCommandSet failed: %v", err)
	}

	if repo.Exists("test") {
		t.Error("Expected 'test' to be deleted")
	}
}

