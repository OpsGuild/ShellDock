package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type ArgumentDef struct {
	Name     string `yaml:"name"`
	Prompt   string `yaml:"prompt,omitempty"`
	Default  string `yaml:"default,omitempty"`
	Required bool   `yaml:"required,omitempty"`
}

type Command struct {
	Description string            `yaml:"description"`
	Command     string            `yaml:"command,omitempty"`
	Platforms   map[string]string `yaml:"platforms,omitempty"`
	SkipOnError bool              `yaml:"skip_on_error,omitempty"`
	Args        []ArgumentDef     `yaml:"args,omitempty"`
}

type CommandSet struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Version     string    `yaml:"version"`
	Tag         string    `yaml:"tag,omitempty"`
	Commands    []Command `yaml:"commands"`
}

type VersionInfo struct {
	Version     string    `yaml:"version"`
	Tag         string    `yaml:"tag,omitempty"`
	Description string    `yaml:"description"`
	Latest      bool      `yaml:"latest,omitempty"`
	Default     bool      `yaml:"default,omitempty"`
	Commands    []Command `yaml:"commands"`
}

type VersionedCommandSet struct {
	Name        string        `yaml:"name"`
	Description string        `yaml:"description,omitempty"`
	Versions    []VersionInfo `yaml:"versions"`
}

type Repository struct {
	path      string
	needsSync bool
}

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
			return r.resolveLatest(&versionedCmdSet)
		}

		for i := range versionedCmdSet.Versions {
			if versionedCmdSet.Versions[i].Tag != "" && strings.EqualFold(versionedCmdSet.Versions[i].Tag, version) {
				v := &versionedCmdSet.Versions[i]
				return &CommandSet{
					Name:        versionedCmdSet.Name,
					Description: v.Description,
					Version:     v.Version,
					Tag:         v.Tag,
					Commands:    v.Commands,
				}, nil
			}
		}

		var defaultMatch *VersionInfo
		var firstMatch *VersionInfo
		for i := range versionedCmdSet.Versions {
			v := &versionedCmdSet.Versions[i]
			if v.Version == version || strings.TrimPrefix(v.Version, "v") == strings.TrimPrefix(version, "v") {
				if v.Default {
					defaultMatch = v
					break
				}
				if firstMatch == nil {
					firstMatch = v
				}
			}
		}

		if defaultMatch != nil {
			return &CommandSet{
				Name:        versionedCmdSet.Name,
				Description: defaultMatch.Description,
				Version:     defaultMatch.Version,
				Tag:         defaultMatch.Tag,
				Commands:    defaultMatch.Commands,
			}, nil
		}
		if firstMatch != nil {
			return &CommandSet{
				Name:        versionedCmdSet.Name,
				Description: firstMatch.Description,
				Version:     firstMatch.Version,
				Tag:         firstMatch.Tag,
				Commands:    firstMatch.Commands,
			}, nil
		}

		return nil, fmt.Errorf("command set '%s' version or tag '%s' not found", name, version)
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

func (r *Repository) resolveLatest(vcs *VersionedCommandSet) (*CommandSet, error) {
	for i := range vcs.Versions {
		if vcs.Versions[i].Latest {
			v := &vcs.Versions[i]
			return &CommandSet{
				Name:        vcs.Name,
				Description: v.Description,
				Version:     v.Version,
				Tag:         v.Tag,
				Commands:    v.Commands,
			}, nil
		}
	}

	highestNum := 0
	var best *VersionInfo
	for i := range vcs.Versions {
		num := extractVersionNumber(vcs.Versions[i].Version)
		if num > highestNum {
			highestNum = num
			best = &vcs.Versions[i]
		}
	}
	if best != nil {
		return &CommandSet{
			Name:        vcs.Name,
			Description: best.Description,
			Version:     best.Version,
			Tag:         best.Tag,
			Commands:    best.Commands,
		}, nil
	}

	if len(vcs.Versions) > 0 {
		v := &vcs.Versions[0]
		return &CommandSet{
			Name:        vcs.Name,
			Description: v.Description,
			Version:     v.Version,
			Tag:         v.Tag,
			Commands:    v.Commands,
		}, nil
	}

	return nil, fmt.Errorf("command set '%s' has no versions", vcs.Name)
}

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

		latestIdx := -1
		for i, v := range versionedCmdSet.Versions {
			if v.Latest {
				latestIdx = i
				break
			}
		}
		if latestIdx == -1 {
			highestNum := 0
			for i, v := range versionedCmdSet.Versions {
				num := extractVersionNumber(v.Version)
				if num > highestNum {
					highestNum = num
					latestIdx = i
				}
			}
		}

		for i, v := range versionedCmdSet.Versions {
			var versionStr string
			if v.Tag != "" {
				versionStr = fmt.Sprintf("%s @%s", v.Version, v.Tag)
			} else {
				versionStr = v.Version
			}
			parts := []string{versionStr}
			if v.Default {
				parts = append(parts, "default")
			}
			if i == latestIdx {
				parts = append(parts, "latest")
			}
			if len(parts) > 1 {
				versions = append(versions, versionStr+" ("+strings.Join(parts[1:], ", ")+")")
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
				t := versionedCmdSet.Versions[i].Tag
				versionMatch := v == versionToSave || strings.TrimPrefix(v, "v") == strings.TrimPrefix(versionToSave, "v")
				tagMatch := cmdSet.Tag != "" && strings.EqualFold(t, cmdSet.Tag)
				if (cmdSet.Tag != "" && versionMatch && tagMatch) || (cmdSet.Tag == "" && versionMatch && t == "") {
					versionedCmdSet.Versions[i].Description = cmdSet.Description
					versionedCmdSet.Versions[i].Tag = cmdSet.Tag
					versionedCmdSet.Versions[i].Commands = cmdSet.Commands
					versionExists = true
					break
				}
			}

			if !versionExists {
				versionedCmdSet.Versions = append(versionedCmdSet.Versions, VersionInfo{
					Version:     versionToSave,
					Tag:         cmdSet.Tag,
					Description: cmdSet.Description,
					Commands:    cmdSet.Commands,
					Latest:      false,
					Default:     false,
				})
			}

			hasLatest := false
			for _, v := range versionedCmdSet.Versions {
				if v.Latest {
					hasLatest = true
					break
				}
			}

			if !hasLatest {
				highestVersionNum := 0
				highestIdx := 0
				for i, v := range versionedCmdSet.Versions {
					versionNum := extractVersionNumber(v.Version)
					if versionNum > highestVersionNum {
						highestVersionNum = versionNum
						highestIdx = i
					}
				}
				versionedCmdSet.Versions[highestIdx].Latest = true
			}

			if versionedCmdSet.Name == "" {
				versionedCmdSet.Name = cmdSet.Name
			}
		} else {
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
						Tag:         oldCmdSet.Tag,
						Description: oldCmdSet.Description,
						Commands:    oldCmdSet.Commands,
						Latest:      oldVersionNum >= newVersionNum,
						Default:     true,
					},
					{
						Version:     versionToSave,
						Tag:         cmdSet.Tag,
						Description: cmdSet.Description,
						Commands:    cmdSet.Commands,
						Latest:      newVersionNum > oldVersionNum,
						Default:     false,
					},
				}
			} else {
				versionedCmdSet.Name = cmdSet.Name
				versionedCmdSet.Versions = []VersionInfo{
					{
						Version:     versionToSave,
						Tag:         cmdSet.Tag,
						Description: cmdSet.Description,
						Commands:    cmdSet.Commands,
						Latest:      true,
						Default:     true,
					},
				}
			}
		}
	} else {
		versionedCmdSet.Name = cmdSet.Name
		versionedCmdSet.Versions = []VersionInfo{
			{
				Version:     versionToSave,
				Tag:         cmdSet.Tag,
				Description: cmdSet.Description,
				Commands:    cmdSet.Commands,
				Latest:      true,
				Default:     true,
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

func (r *Repository) FindCommandSetFile(name string) string {
	return r.findCommandSetFile(name)
}

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

type CommandGroup struct {
	Name     string
	Commands []string
}

func (r *Repository) ListCommandSetsGrouped() ([]CommandGroup, error) {
	if _, err := os.Stat(r.path); os.IsNotExist(err) {
		return []CommandGroup{}, nil
	}

	groupMap := make(map[string][]string)
	var groupOrder []string

	err := filepath.WalkDir(r.path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && filepath.Ext(d.Name()) == ".yaml" {
			name := d.Name()[:len(d.Name())-5]

			relPath, relErr := filepath.Rel(r.path, path)
			group := "general"
			if relErr == nil {
				dir := filepath.Dir(relPath)
				if dir != "." {
					parts := strings.Split(filepath.ToSlash(dir), "/")
					group = parts[0]
				}
			}

			if _, exists := groupMap[group]; !exists {
				groupOrder = append(groupOrder, group)
			}
			groupMap[group] = append(groupMap[group], name)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to read repository: %w", err)
	}

	var groups []CommandGroup
	for _, groupName := range groupOrder {
		groups = append(groups, CommandGroup{
			Name:     groupName,
			Commands: groupMap[groupName],
		})
	}

	return groups, nil
}

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

func (r *Repository) Exists(name string) bool {
	return r.findCommandSetFile(name) != ""
}

func (r *Repository) GetPath() string {
	return r.path
}

func (r *Repository) NeedsSync() bool {
	return r.needsSync
}
