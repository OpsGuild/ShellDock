package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}

	if path == "" {
		t.Error("Expected non-empty config path")
	}

	expectedSuffix := filepath.Join(".shelldock", ConfigFileName)
	if !strings.HasSuffix(path, expectedSuffix) {
		t.Errorf("Expected path to end with %s, got %s", expectedSuffix, path)
	}
}

func TestLoadConfig_Default(t *testing.T) {
	// Test with non-existent config (should return default)
	// We can't easily mock GetConfigPath, so we test the actual behavior
	// by using a temp directory that doesn't have a config file
	tmpDir := t.TempDir()
	
	// Temporarily override HOME
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			_ = os.Setenv("HOME", originalHome)
		}
	}()
	
	_ = os.Setenv("HOME", tmpDir)
	
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Platform != DefaultPlatform {
		t.Errorf("Expected platform '%s', got '%s'", DefaultPlatform, cfg.Platform)
	}
}

func TestLoadConfig_FromFile(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Temporarily override HOME
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			_ = os.Setenv("HOME", originalHome)
		}
	}()
	
	_ = os.Setenv("HOME", tmpDir)
	
	// Create config directory and file
	configDir := filepath.Join(tmpDir, ".shelldock")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}
	configPath := filepath.Join(configDir, ConfigFileName)
	
	// Create config file
	configContent := "platform: ubuntu\n"
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Platform != "ubuntu" {
		t.Errorf("Expected platform 'ubuntu', got '%s'", cfg.Platform)
	}
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Temporarily override HOME
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			_ = os.Setenv("HOME", originalHome)
		}
	}()
	
	_ = os.Setenv("HOME", tmpDir)

	cfg := &Config{Platform: "centos"}
	err := SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify it was saved
	saved, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if saved.Platform != "centos" {
		t.Errorf("Expected platform 'centos', got '%s'", saved.Platform)
	}
}

func TestGetPlatform_Auto(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Temporarily override HOME
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			_ = os.Setenv("HOME", originalHome)
		}
	}()
	
	_ = os.Setenv("HOME", tmpDir)

	cfg := &Config{Platform: "auto"}
	err := SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	platform, err := GetPlatform()
	if err != nil {
		t.Fatalf("GetPlatform failed: %v", err)
	}

	if platform == "" {
		t.Error("Expected non-empty platform")
	}
}

func TestDetectPlatform(t *testing.T) {
	platform := DetectPlatform()
	if platform == "" {
		t.Error("Expected non-empty platform")
	}

	// Should be one of the known platforms
	validPlatforms := []string{"linux", "darwin", "windows", "ubuntu", "debian", "centos", "fedora", "arch"}
	valid := false
	for _, vp := range validPlatforms {
		if platform == vp {
			valid = true
			break
		}
	}
	if !valid {
		t.Errorf("Platform '%s' is not a recognized platform", platform)
	}
}

func TestNormalizeDistroID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ubuntu", "ubuntu"},
		{"Ubuntu", "ubuntu"},
		{"UBUNTU", "ubuntu"},
		{"debian", "debian"},
		{"centos", "centos"},
		{"rhel", "rhel"},
		{"redhat", "rhel"},
		{"fedora", "fedora"},
		{"arch", "arch"},
		{"archlinux", "arch"},
		{"unknown", ""},
	}

	for _, tt := range tests {
		result := normalizeDistroID(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeDistroID(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

