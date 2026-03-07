package repo

import (
	"os"
	"path/filepath"
	"strings"
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
    latest: false
    default: true
    description: Version 1
    commands:
      - description: Command 1
        command: echo "v1"
  - version: "v2"
    latest: true
    default: true
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

	cmdSet, err := repo.GetCommandSet("test", "")
	if err != nil {
		t.Fatalf("GetCommandSet failed: %v", err)
	}
	if cmdSet.Version != "v2" {
		t.Errorf("Expected latest version 'v2', got %s", cmdSet.Version)
	}

	cmdSet, err = repo.GetCommandSet("test", "v1")
	if err != nil {
		t.Fatalf("GetCommandSet failed: %v", err)
	}
	if cmdSet.Version != "v1" {
		t.Errorf("Expected version 'v1', got %s", cmdSet.Version)
	}
}

func TestGetCommandSet_VersionedWithTags(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewRepository(tmpDir)

	yamlContent := `name: mysql
description: MySQL/MariaDB
versions:
  - version: "v1"
    tag: mysql
    latest: true
    default: true
    description: MySQL server
    commands:
      - description: Install MySQL
        command: echo "mysql"
  - version: "v1"
    tag: mariadb
    latest: false
    default: false
    description: MariaDB server
    commands:
      - description: Install MariaDB
        command: echo "mariadb"
`
	filePath := filepath.Join(tmpDir, "mysql.yaml")
	err := os.WriteFile(filePath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cmdSet, err := repo.GetCommandSet("mysql", "")
	if err != nil {
		t.Fatalf("GetCommandSet failed: %v", err)
	}
	if cmdSet.Tag != "mysql" {
		t.Errorf("Expected latest to resolve to tag 'mysql', got '%s'", cmdSet.Tag)
	}

	cmdSet, err = repo.GetCommandSet("mysql", "mariadb")
	if err != nil {
		t.Fatalf("GetCommandSet by tag failed: %v", err)
	}
	if cmdSet.Tag != "mariadb" {
		t.Errorf("Expected tag 'mariadb', got '%s'", cmdSet.Tag)
	}

	cmdSet, err = repo.GetCommandSet("mysql", "mysql")
	if err != nil {
		t.Fatalf("GetCommandSet by tag failed: %v", err)
	}
	if cmdSet.Tag != "mysql" {
		t.Errorf("Expected tag 'mysql', got '%s'", cmdSet.Tag)
	}

	cmdSet, err = repo.GetCommandSet("mysql", "v1")
	if err != nil {
		t.Fatalf("GetCommandSet by version failed: %v", err)
	}
	if cmdSet.Tag != "mysql" {
		t.Errorf("Expected default for v1 to be tag 'mysql', got '%s'", cmdSet.Tag)
	}
}

func TestGetCommandSet_VersionedDefaultDisambiguation(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewRepository(tmpDir)

	yamlContent := `name: certbot
description: Certbot SSL
versions:
  - version: "v1"
    tag: certonly
    latest: false
    default: false
    description: Standalone
    commands:
      - description: Install
        command: echo "certonly"
  - version: "v1"
    tag: nginx
    latest: true
    default: true
    description: With Nginx
    commands:
      - description: Install
        command: echo "nginx"
`
	filePath := filepath.Join(tmpDir, "certbot.yaml")
	err := os.WriteFile(filePath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cmdSet, err := repo.GetCommandSet("certbot", "v1")
	if err != nil {
		t.Fatalf("GetCommandSet by version failed: %v", err)
	}
	if cmdSet.Tag != "nginx" {
		t.Errorf("Expected default for v1 to be tag 'nginx', got '%s'", cmdSet.Tag)
	}

	cmdSet, err = repo.GetCommandSet("certbot", "certonly")
	if err != nil {
		t.Fatalf("GetCommandSet by tag failed: %v", err)
	}
	if cmdSet.Tag != "certonly" {
		t.Errorf("Expected tag 'certonly', got '%s'", cmdSet.Tag)
	}

	cmdSet, err = repo.GetCommandSet("certbot", "")
	if err != nil {
		t.Fatalf("GetCommandSet latest failed: %v", err)
	}
	if cmdSet.Tag != "nginx" {
		t.Errorf("Expected latest to be 'nginx', got '%s'", cmdSet.Tag)
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
    latest: false
    default: true
    commands: []
  - version: "v2"
    latest: true
    default: true
    commands: []
  - version: "v3"
    latest: false
    default: true
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

	foundLatest := false
	for _, v := range versions {
		if strings.Contains(v, "v2") && strings.Contains(v, "latest") {
			foundLatest = true
			break
		}
	}
	if !foundLatest {
		t.Errorf("Expected to find v2 marked as latest in versions, got %v", versions)
	}
}

func TestListVersions_WithTags(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewRepository(tmpDir)

	yamlContent := `name: mysql
versions:
  - version: "v1"
    tag: mysql
    latest: true
    default: true
    commands: []
  - version: "v1"
    tag: mariadb
    latest: false
    default: false
    commands: []
`
	filePath := filepath.Join(tmpDir, "mysql.yaml")
	err := os.WriteFile(filePath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	versions, err := repo.ListVersions("mysql")
	if err != nil {
		t.Fatalf("ListVersions failed: %v", err)
	}

	if len(versions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(versions))
	}

	foundDefault := false
	foundLatest := false
	for _, v := range versions {
		if strings.Contains(v, "@mysql") && strings.Contains(v, "default") {
			foundDefault = true
		}
		if strings.Contains(v, "@mysql") && strings.Contains(v, "latest") {
			foundLatest = true
		}
	}
	if !foundDefault {
		t.Errorf("Expected mysql entry marked as default, got %v", versions)
	}
	if !foundLatest {
		t.Errorf("Expected mysql entry marked as latest, got %v", versions)
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
