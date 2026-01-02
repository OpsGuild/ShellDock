package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// ArgumentDef represents a command argument definition
type ArgumentDef struct {
	Name     string `yaml:"name"`               // Variable name (e.g., "username")
	Prompt   string `yaml:"prompt,omitempty"`   // Prompt question (e.g., "Enter your name:")
	Default  string `yaml:"default,omitempty"`  // Default value
	Required bool   `yaml:"required,omitempty"` // Whether argument is required
}

// Command represents a single command step
type Command struct {
	Description string            `yaml:"description"`
	Command     string            `yaml:"command,omitempty"`   // Single command (backward compatibility)
	Platforms   map[string]string `yaml:"platforms,omitempty"` // Platform-specific commands: platform -> command
	SkipOnError bool              `yaml:"skip_on_error,omitempty"`
	Args        []ArgumentDef     `yaml:"args,omitempty"` // Argument definitions for this command
}

// CommandSet represents a collection of commands for a topic
type CommandSet struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Version     string    `yaml:"version"`
	Commands    []Command `yaml:"commands"`
}

// VersionInfo represents a single version of a command set
type VersionInfo struct {
	Version     string    `yaml:"version"`
	Description string    `yaml:"description"`
	Latest      bool      `yaml:"latest,omitempty"` // Mark this version as latest
	Commands    []Command `yaml:"commands"`
}

// VersionedCommandSet represents a command set with multiple versions
type VersionedCommandSet struct {
	Name        string        `yaml:"name"`
	Description string        `yaml:"description,omitempty"`
	Versions    []VersionInfo `yaml:"versions"` // Array of versions
}

// Repository manages command sets
type Repository struct {
	path string
}

// NewRepository creates a new repository instance
func NewRepository(path string) *Repository {
	return &Repository{path: path}
}

// extractVersionNumber extracts numeric version from "v1", "v2", etc.
func extractVersionNumber(version string) int {
	version = strings.TrimPrefix(strings.ToLower(version), "v")
	num, err := strconv.Atoi(version)
	if err != nil {
		return 0
	}
	return num
}

// GetCommandSet retrieves a command set by name and optional version
// If version is empty, returns the latest version
func (r *Repository) GetCommandSet(name string, version string) (*CommandSet, error) {
	filePath := filepath.Join(r.path, fmt.Sprintf("%s.yaml", name))

	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("command set '%s' not found", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read command set: %w", err)
	}

	// Try to parse as versioned command set first
	var versionedCmdSet VersionedCommandSet
	if err := yaml.Unmarshal(data, &versionedCmdSet); err == nil && versionedCmdSet.Versions != nil && len(versionedCmdSet.Versions) > 0 {
		// It's a versioned command set
		if version == "" || version == "latest" {
			// Find the latest version (marked with latest: true or highest version number)
			latestVersion := ""
			hasLatestFlag := false
			highestVersionNum := 0

			// First, check for latest flag
			for _, v := range versionedCmdSet.Versions {
				if v.Latest {
					version = v.Version
					hasLatestFlag = true
					break
				}
			}

			// If no latest flag, find highest version number
			if !hasLatestFlag {
				for _, v := range versionedCmdSet.Versions {
					versionNum := extractVersionNumber(v.Version)
					if versionNum > highestVersionNum {
						highestVersionNum = versionNum
						latestVersion = v.Version
					}
				}
				if latestVersion != "" {
					version = latestVersion
				}
			}
		}

		// Find the requested version
		var foundVersion *VersionInfo
		for i := range versionedCmdSet.Versions {
			// Support both "v1" and "1" formats
			v := versionedCmdSet.Versions[i].Version
			if v == version || strings.TrimPrefix(v, "v") == strings.TrimPrefix(version, "v") {
				foundVersion = &versionedCmdSet.Versions[i]
				break
			}
		}

		if foundVersion == nil {
			return nil, fmt.Errorf("command set '%s' version '%s' not found", name, version)
		}

		// Convert VersionInfo to CommandSet
		cmdSet := CommandSet{
			Name:        versionedCmdSet.Name,
			Description: foundVersion.Description,
			Version:     foundVersion.Version,
			Commands:    foundVersion.Commands,
		}

		return &cmdSet, nil
	}

	// Fallback to single version format (backward compatibility)
	var cmdSet CommandSet
	if err := yaml.Unmarshal(data, &cmdSet); err != nil {
		return nil, fmt.Errorf("failed to parse command set: %w", err)
	}

	// If version was specified but file is single-version format, check if it matches
	if version != "" && version != "latest" {
		// Support both "v1" and "1" formats
		if cmdSet.Version != version && strings.TrimPrefix(cmdSet.Version, "v") != strings.TrimPrefix(version, "v") {
			return nil, fmt.Errorf("command set '%s' version '%s' not found (file contains version '%s')", name, version, cmdSet.Version)
		}
	}

	return &cmdSet, nil
}

// ListVersions returns all available versions for a command set
func (r *Repository) ListVersions(name string) ([]string, error) {
	filePath := filepath.Join(r.path, fmt.Sprintf("%s.yaml", name))

	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return []string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read command set: %w", err)
	}

	// Try to parse as versioned command set
	var versionedCmdSet VersionedCommandSet
	if err := yaml.Unmarshal(data, &versionedCmdSet); err == nil && versionedCmdSet.Versions != nil && len(versionedCmdSet.Versions) > 0 {
		// It's a versioned command set
		var versions []string
		var latestVersion string

		// Find latest version
		for _, v := range versionedCmdSet.Versions {
			if v.Latest {
				latestVersion = v.Version
				break
			}
		}

		// If no latest flag, find highest version number
		if latestVersion == "" {
			highestVersionNum := 0
			for _, v := range versionedCmdSet.Versions {
				versionNum := extractVersionNumber(v.Version)
				if versionNum > highestVersionNum {
					highestVersionNum = versionNum
					latestVersion = v.Version
				}
			}
		}

		// Build version list
		for _, v := range versionedCmdSet.Versions {
			if v.Version == latestVersion {
				versions = append(versions, v.Version+" (latest)")
			} else {
				versions = append(versions, v.Version)
			}
		}

		return versions, nil
	}

	// Fallback: single version format
	var cmdSet CommandSet
	if err := yaml.Unmarshal(data, &cmdSet); err == nil {
		version := cmdSet.Version
		if version == "" {
			version = "v1"
		}
		return []string{version + " (latest)"}, nil
	}

	return []string{}, nil
}

// SaveCommandSet saves a command set to the repository
// If version is specified, adds/updates that version in the versioned file
// If version is empty, uses the version from cmdSet.Version
func (r *Repository) SaveCommandSet(cmdSet *CommandSet, version string) error {
	if err := os.MkdirAll(r.path, 0755); err != nil {
		return fmt.Errorf("failed to create repository directory: %w", err)
	}

	filePath := filepath.Join(r.path, fmt.Sprintf("%s.yaml", cmdSet.Name))

	// Determine which version to save
	versionToSave := version
	if versionToSave == "" || versionToSave == "latest" {
		versionToSave = cmdSet.Version
		if versionToSave == "" {
			versionToSave = "v1" // Default version
		}
	}

	// Ensure version starts with 'v' if it's numeric
	if !strings.HasPrefix(versionToSave, "v") {
		if _, err := strconv.Atoi(versionToSave); err == nil {
			versionToSave = "v" + versionToSave
		}
	}

	// Try to load existing versioned command set
	var versionedCmdSet VersionedCommandSet
	data, err := os.ReadFile(filePath)
	if err == nil {
		// File exists, try to parse as versioned
		if err := yaml.Unmarshal(data, &versionedCmdSet); err == nil && versionedCmdSet.Versions != nil && len(versionedCmdSet.Versions) > 0 {
			// It's already a versioned command set
			// Check if version already exists
			versionExists := false
			for i := range versionedCmdSet.Versions {
				v := versionedCmdSet.Versions[i].Version
				if v == versionToSave || strings.TrimPrefix(v, "v") == strings.TrimPrefix(versionToSave, "v") {
					// Update existing version
					versionedCmdSet.Versions[i].Description = cmdSet.Description
					versionedCmdSet.Versions[i].Commands = cmdSet.Commands
					versionExists = true
					break
				}
			}

			// Add new version if it doesn't exist
			if !versionExists {
				versionedCmdSet.Versions = append(versionedCmdSet.Versions, VersionInfo{
					Version:     versionToSave,
					Description: cmdSet.Description,
					Commands:    cmdSet.Commands,
					Latest:      false, // Will be set below if needed
				})
			}

			// Find highest version number to mark as latest
			highestVersionNum := 0
			latestVersion := ""
			for _, v := range versionedCmdSet.Versions {
				versionNum := extractVersionNumber(v.Version)
				if versionNum > highestVersionNum {
					highestVersionNum = versionNum
					latestVersion = v.Version
				}
			}

			// Update latest flags
			for i := range versionedCmdSet.Versions {
				versionedCmdSet.Versions[i].Latest = (versionedCmdSet.Versions[i].Version == latestVersion)
			}

			if versionedCmdSet.Name == "" {
				versionedCmdSet.Name = cmdSet.Name
			}
		} else {
			// Existing file is single-version format, convert it
			var oldCmdSet CommandSet
			if err := yaml.Unmarshal(data, &oldCmdSet); err == nil {
				versionedCmdSet.Name = oldCmdSet.Name
				oldVersion := oldCmdSet.Version
				if oldVersion == "" {
					oldVersion = "v1"
				}
				// Ensure old version has 'v' prefix
				if !strings.HasPrefix(oldVersion, "v") {
					if _, err := strconv.Atoi(oldVersion); err == nil {
						oldVersion = "v" + oldVersion
					}
				}

				// Determine which is latest
				oldVersionNum := extractVersionNumber(oldVersion)
				newVersionNum := extractVersionNumber(versionToSave)

				// Create versions array
				versionedCmdSet.Versions = []VersionInfo{
					{
						Version:     oldVersion,
						Description: oldCmdSet.Description,
						Commands:    oldCmdSet.Commands,
						Latest:      oldVersionNum >= newVersionNum,
					},
					{
						Version:     versionToSave,
						Description: cmdSet.Description,
						Commands:    cmdSet.Commands,
						Latest:      newVersionNum > oldVersionNum,
					},
				}
			} else {
				// Can't parse, create new
				versionedCmdSet.Name = cmdSet.Name
				versionedCmdSet.Versions = []VersionInfo{
					{
						Version:     versionToSave,
						Description: cmdSet.Description,
						Commands:    cmdSet.Commands,
						Latest:      true,
					},
				}
			}
		}
	} else {
		// File doesn't exist, create new versioned command set
		versionedCmdSet.Name = cmdSet.Name
		versionedCmdSet.Versions = []VersionInfo{
			{
				Version:     versionToSave,
				Description: cmdSet.Description,
				Commands:    cmdSet.Commands,
				Latest:      true,
			},
		}
	}

	// Marshal and save
	data, err = yaml.Marshal(&versionedCmdSet)
	if err != nil {
		return fmt.Errorf("failed to marshal command set: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write command set: %w", err)
	}

	return nil
}

// ListCommandSets returns all available command sets
func (r *Repository) ListCommandSets() ([]string, error) {
	if _, err := os.Stat(r.path); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(r.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read repository: %w", err)
	}

	var sets []string

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".yaml" {
			name := entry.Name()[:len(entry.Name())-5] // Remove .yaml
			sets = append(sets, name)
		}
	}

	return sets, nil
}

// DeleteCommandSet removes a command set from the repository
func (r *Repository) DeleteCommandSet(name string) error {
	filePath := filepath.Join(r.path, fmt.Sprintf("%s.yaml", name))

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("command set '%s' not found", name)
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete command set: %w", err)
	}

	return nil
}

// Exists checks if a command set exists
func (r *Repository) Exists(name string) bool {
	filePath := filepath.Join(r.path, fmt.Sprintf("%s.yaml", name))
	_, err := os.Stat(filePath)
	return err == nil
}
