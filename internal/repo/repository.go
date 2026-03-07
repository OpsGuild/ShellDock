package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// ArgumentDef represents a command argument definition.
type ArgumentDef struct {
	Name     string `yaml:"name"`
	Prompt   string `yaml:"prompt,omitempty"`
	Default  string `yaml:"default,omitempty"`
	Required bool   `yaml:"required,omitempty"`
}

// Command represents a single command step.
type Command struct {
	Description string            `yaml:"description"`
	Command     string            `yaml:"command,omitempty"`
	Platforms   map[string]string `yaml:"platforms,omitempty"`
	SkipOnError bool              `yaml:"skip_on_error,omitempty"`
	Args        []ArgumentDef     `yaml:"args,omitempty"`
}

// CommandSet represents a collection of commands for a topic.
type CommandSet struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Version     string    `yaml:"version"`
	Commands    []Command `yaml:"commands"`
}

// VersionInfo represents a single version of a command set.
type VersionInfo struct {
	Version     string    `yaml:"version"`
	Tag         string    `yaml:"tag,omitempty"`
	Description string    `yaml:"description"`
	Latest      bool      `yaml:"latest,omitempty"`
	Commands    []Command `yaml:"commands"`
}

// VersionedCommandSet represents a command set with multiple versions.
type VersionedCommandSet struct {
	Name        string        `yaml:"name"`
	Description string        `yaml:"description,omitempty"`
	Versions    []VersionInfo `yaml:"versions"`
}

// Repository manages command sets.
type Repository struct {
	path      string
	needsSync bool
}

// NewRepository creates a new repository instance.
func NewRepository(path string) *Repository {
	return &Repository{path: path}
}

func extractVersionNumber(version string) int {
	version = strings.TrimPrefix(strings.ToLower(version), "v")
	num, err := strconv.Atoi(version)
	if err != nil {
		return 0
	}
	return num
}

// GetCommandSet retrieves a command set by name and optional version.
// If version is empty, returns the latest version.
func (r *Repository) GetCommandSet(name string, version string) (*CommandSet, error) {
	filePath := r.findCommandSetFile(name)
	if filePath == "" {
		return nil, fmt.Errorf("command set '%s' not found", name)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read command set: %w", err)
	}

	var versionedCmdSet VersionedCommandSet
	if err := yaml.Unmarshal(data, &versionedCmdSet); err == nil && versionedCmdSet.Versions != nil && len(versionedCmdSet.Versions) > 0 {
		if version == "" || version == "latest" {
			latestVersion := ""
			hasLatestFlag := false
			highestVersionNum := 0

			for _, v := range versionedCmdSet.Versions {
				if v.Latest {
					version = v.Version
					hasLatestFlag = true
					break
				}
			}

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

		var foundVersion *VersionInfo
		for i := range versionedCmdSet.Versions {
			v := versionedCmdSet.Versions[i]
			if v.Version == version || strings.TrimPrefix(v.Version, "v") == strings.TrimPrefix(version, "v") {
				foundVersion = &versionedCmdSet.Versions[i]
				break
			}
			if v.Tag != "" && strings.EqualFold(v.Tag, version) {
				foundVersion = &versionedCmdSet.Versions[i]
				break
			}
		}

		if foundVersion == nil {
			return nil, fmt.Errorf("command set '%s' version or tag '%s' not found", name, version)
		}

		cmdSet := CommandSet{
			Name:        versionedCmdSet.Name,
			Description: foundVersion.Description,
			Version:     foundVersion.Version,
			Commands:    foundVersion.Commands,
		}

		return &cmdSet, nil
	}

	var cmdSet CommandSet
	if err := yaml.Unmarshal(data, &cmdSet); err != nil {
		return nil, fmt.Errorf("failed to parse command set: %w", err)
	}

	if version != "" && version != "latest" {
		if cmdSet.Version != version && strings.TrimPrefix(cmdSet.Version, "v") != strings.TrimPrefix(version, "v") {
			return nil, fmt.Errorf("command set '%s' version '%s' not found (file contains version '%s')", name, version, cmdSet.Version)
		}
	}

	return &cmdSet, nil
}

// ListVersions returns all available versions for a command set.
func (r *Repository) ListVersions(name string) ([]string, error) {
	filePath := r.findCommandSetFile(name)
	if filePath == "" {
		return []string{}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read command set: %w", err)
	}

	var versionedCmdSet VersionedCommandSet
	if err := yaml.Unmarshal(data, &versionedCmdSet); err == nil && versionedCmdSet.Versions != nil && len(versionedCmdSet.Versions) > 0 {
		var versions []string
		var latestVersion string

		for _, v := range versionedCmdSet.Versions {
			if v.Latest {
				latestVersion = v.Version
				break
			}
		}

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

		for _, v := range versionedCmdSet.Versions {
			versionStr := v.Version
			if v.Tag != "" {
				versionStr = fmt.Sprintf("%s [%s]", v.Version, v.Tag)
			}
			if v.Version == latestVersion {
				versions = append(versions, versionStr+" (latest)")
			} else {
				versions = append(versions, versionStr)
			}
		}

		return versions, nil
	}

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

// SaveCommandSet saves a command set to the repository.
// If version is specified, adds/updates that version in the versioned file.
// If version is empty, uses the version from cmdSet.Version.
func (r *Repository) SaveCommandSet(cmdSet *CommandSet, version string) error {
	if err := os.MkdirAll(r.path, 0755); err != nil {
		return fmt.Errorf("failed to create repository directory: %w", err)
	}

	filePath := filepath.Join(r.path, fmt.Sprintf("%s.yaml", cmdSet.Name))

	versionToSave := version
	if versionToSave == "" || versionToSave == "latest" {
		versionToSave = cmdSet.Version
		if versionToSave == "" {
			versionToSave = "v1"
		}
	}

	if !strings.HasPrefix(versionToSave, "v") {
		if _, err := strconv.Atoi(versionToSave); err == nil {
			versionToSave = "v" + versionToSave
		}
	}

	var versionedCmdSet VersionedCommandSet
	data, err := os.ReadFile(filePath)
	if err == nil {
		if err := yaml.Unmarshal(data, &versionedCmdSet); err == nil && versionedCmdSet.Versions != nil && len(versionedCmdSet.Versions) > 0 {
			versionExists := false
			for i := range versionedCmdSet.Versions {
				v := versionedCmdSet.Versions[i].Version
				if v == versionToSave || strings.TrimPrefix(v, "v") == strings.TrimPrefix(versionToSave, "v") {
					versionedCmdSet.Versions[i].Description = cmdSet.Description
					versionedCmdSet.Versions[i].Commands = cmdSet.Commands
					versionExists = true
					break
				}
			}

			if !versionExists {
				versionedCmdSet.Versions = append(versionedCmdSet.Versions, VersionInfo{
					Version:     versionToSave,
					Description: cmdSet.Description,
					Commands:    cmdSet.Commands,
					Latest:      false,
				})
			}

			highestVersionNum := 0
			latestVersion := ""
			for _, v := range versionedCmdSet.Versions {
				versionNum := extractVersionNumber(v.Version)
				if versionNum > highestVersionNum {
					highestVersionNum = versionNum
					latestVersion = v.Version
				}
			}

			for i := range versionedCmdSet.Versions {
				versionedCmdSet.Versions[i].Latest = (versionedCmdSet.Versions[i].Version == latestVersion)
			}

			if versionedCmdSet.Name == "" {
				versionedCmdSet.Name = cmdSet.Name
			}
		} else {
			// Existing file is single-version format, convert to versioned
			var oldCmdSet CommandSet
			if err := yaml.Unmarshal(data, &oldCmdSet); err == nil {
				versionedCmdSet.Name = oldCmdSet.Name
				oldVersion := oldCmdSet.Version
				if oldVersion == "" {
					oldVersion = "v1"
				}
				if !strings.HasPrefix(oldVersion, "v") {
					if _, err := strconv.Atoi(oldVersion); err == nil {
						oldVersion = "v" + oldVersion
					}
				}

				oldVersionNum := extractVersionNumber(oldVersion)
				newVersionNum := extractVersionNumber(versionToSave)

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

	data, err = yaml.Marshal(&versionedCmdSet)
	if err != nil {
		return fmt.Errorf("failed to marshal command set: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write command set: %w", err)
	}

	return nil
}

func (r *Repository) findCommandSetFile(name string) string {
	rootPath := filepath.Join(r.path, fmt.Sprintf("%s.yaml", name))
	if _, err := os.Stat(rootPath); err == nil {
		return rootPath
	}

	var foundPath string
	_ = filepath.WalkDir(r.path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && filepath.Ext(d.Name()) == ".yaml" {
			baseName := d.Name()[:len(d.Name())-5]
			if baseName == name {
				foundPath = path
				return filepath.SkipAll
			}
		}
		return nil
	})

	return foundPath
}

// ListCommandSets returns all available command sets.
func (r *Repository) ListCommandSets() ([]string, error) {
	if _, err := os.Stat(r.path); os.IsNotExist(err) {
		return []string{}, nil
	}

	var sets []string

	err := filepath.WalkDir(r.path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && filepath.Ext(d.Name()) == ".yaml" {
			name := d.Name()[:len(d.Name())-5]
			sets = append(sets, name)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to read repository: %w", err)
	}

	return sets, nil
}

// DeleteCommandSet removes a command set from the repository.
func (r *Repository) DeleteCommandSet(name string) error {
	filePath := r.findCommandSetFile(name)
	if filePath == "" {
		return fmt.Errorf("command set '%s' not found", name)
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete command set: %w", err)
	}

	return nil
}

// Exists checks if a command set exists.
func (r *Repository) Exists(name string) bool {
	return r.findCommandSetFile(name) != ""
}

// GetPath returns the repository path.
func (r *Repository) GetPath() string {
	return r.path
}

// NeedsSync returns true if the repository needs initial sync.
func (r *Repository) NeedsSync() bool {
	return r.needsSync
}
