package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	ConfigFileName  = ".sdrc"
	DefaultPlatform = "auto"
)

// Config represents the ShellDock configuration.
type Config struct {
	Platform string `yaml:"platform"`
}

// GetConfigPath returns the path to the config file.
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".shelldock", ConfigFileName), nil
}

// LoadConfig loads the configuration from ~/.shelldock/.sdrc
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	config := &Config{
		Platform: DefaultPlatform,
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// SaveConfig saves the configuration to ~/.shelldock/.sdrc
func SaveConfig(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetPlatform returns the current platform.
// If config has "auto", detects the platform automatically.
func GetPlatform() (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", err
	}

	if config.Platform == "" || config.Platform == "auto" {
		return DetectPlatform(), nil
	}

	return config.Platform, nil
}

// DetectPlatform detects the current platform.
// For Linux, it also detects the distribution.
func DetectPlatform() string {
	os := runtime.GOOS
	switch os {
	case "linux":
		return DetectLinuxDistribution()
	case "darwin":
		return "darwin"
	case "windows":
		return "windows"
	default:
		return "linux"
	}
}

// DetectLinuxDistribution detects the Linux distribution.
func DetectLinuxDistribution() string {
	if distro := readOSRelease(); distro != "" {
		return distro
	}

	if _, err := os.Stat("/etc/debian_version"); err == nil {
		return "debian"
	}
	if _, err := os.Stat("/etc/redhat-release"); err == nil {
		return "centos"
	}
	if _, err := os.Stat("/etc/arch-release"); err == nil {
		return "arch"
	}
	if _, err := os.Stat("/etc/SuSE-release"); err == nil {
		return "opensuse"
	}

	return "linux"
}

func readOSRelease() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return ""
	}

	lines := strings.Split(string(data), "\n")
	var id string
	var idLike string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "ID=") {
			id = strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
		}
		if strings.HasPrefix(line, "ID_LIKE=") {
			idLike = strings.Trim(strings.TrimPrefix(line, "ID_LIKE="), "\"")
		}
	}

	if id != "" {
		normalized := normalizeDistroID(id)
		if normalized != "" {
			return normalized
		}
	}

	if idLike != "" {
		normalized := normalizeDistroID(idLike)
		if normalized != "" {
			return normalized
		}
	}

	return ""
}

func normalizeDistroID(id string) string {
	id = strings.ToLower(id)

	distroMap := map[string]string{
		"ubuntu":              "ubuntu",
		"debian":              "debian",
		"centos":              "centos",
		"rhel":                "rhel",
		"redhat":              "rhel",
		"fedora":              "fedora",
		"arch":                "arch",
		"archlinux":           "arch",
		"opensuse":            "opensuse",
		"suse":                "opensuse",
		"opensuse-leap":       "opensuse",
		"opensuse-tumbleweed": "opensuse",
		"alpine":              "alpine",
		"amazon":              "amazon",
		"oracle":              "oracle",
	}

	if normalized, exists := distroMap[id]; exists {
		return normalized
	}

	for key, value := range distroMap {
		if strings.Contains(id, key) {
			return value
		}
	}

	return ""
}
