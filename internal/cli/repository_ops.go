package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/shelldock/shelldock/internal/repo"
	"github.com/spf13/cobra"
)

var (
	addOpen        bool
	addDescription string
	addVersion     string

	rmYes bool

	openEditor string
)

type editorCommand struct {
	bin  string
	args []string
}

var addCmd = &cobra.Command{
	Use:   "add [command-set-name]",
	Short: "Add a command set to local repository",
	Long: `Add a command set to your local repository (~/.shelldock).
If a bundled command set with the same name exists, it will be copied locally for customization.
Otherwise, a new template command set is created.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, err := normalizeCommandSetName(args[0])
		handleError(err)

		manager, err := repo.NewManager()
		handleError(err)

		localRepo := manager.GetLocalRepo()
		if localRepo.Exists(name) {
			fmt.Printf("❌ Command set '%s' already exists in local repository.\n", name)
			fmt.Printf("💡 Use 'shelldock open %s' to edit it.\n", name)
			return
		}

		bundledRepo := manager.GetBundledRepo()
		if bundledRepo != nil && bundledRepo.Exists(name) {
			localPath, err := copyBundledCommandSetToLocal(localRepo, bundledRepo, name)
			if err != nil {
				fmt.Printf("❌ Failed to add local copy of '%s': %v\n", name, err)
				return
			}

			fmt.Printf("✅ Added '%s' to local repository from bundled command set.\n", name)
			fmt.Printf("📄 File: %s\n", localPath)

			if addOpen {
				if err := openInEditor(localPath, openEditor); err != nil {
					fmt.Printf("⚠️  Added command set but failed to open editor: %v\n", err)
				}
			}
			return
		}

		description := addDescription
		if strings.TrimSpace(description) == "" {
			description = fmt.Sprintf("Local command set for %s", name)
		}

		version := strings.TrimSpace(addVersion)
		if version == "" {
			version = "v1"
		}

		cmdSet := &repo.CommandSet{
			Name:        name,
			Description: description,
			Version:     version,
			Commands: []repo.Command{
				{
					Description: "Example step",
					Command:     "echo \"Edit this command\"",
				},
			},
		}

		if err := localRepo.SaveCommandSet(cmdSet, version); err != nil {
			fmt.Printf("❌ Failed to create command set '%s': %v\n", name, err)
			return
		}

		localPath := localRepo.FindCommandSetFile(name)
		fmt.Printf("✅ Created local command set '%s'.\n", name)
		fmt.Printf("📄 File: %s\n", localPath)
		fmt.Printf("💡 Edit with: shelldock open %s\n", name)

		if addOpen {
			if err := openInEditor(localPath, openEditor); err != nil {
				fmt.Printf("⚠️  Created command set but failed to open editor: %v\n", err)
			}
		}
	},
}

var rmCmd = &cobra.Command{
	Use:     "rm [command-set-name]",
	Aliases: []string{"remove", "delete"},
	Short:   "Remove a command set from local repository",
	Long:    "Remove a command set from local repository (~/.shelldock). Bundled repository files are never deleted.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, err := normalizeCommandSetName(args[0])
		handleError(err)

		manager, err := repo.NewManager()
		handleError(err)

		localRepo := manager.GetLocalRepo()
		if !localRepo.Exists(name) {
			if bundled := manager.GetBundledRepo(); bundled != nil && bundled.Exists(name) {
				fmt.Printf("❌ Command set '%s' exists only in bundled repository and cannot be removed with 'rm'.\n", name)
				fmt.Println("💡 Copy it locally first with 'shelldock add <name>' if you want a local editable version.")
				return
			}
			fmt.Printf("❌ Command set '%s' not found in local repository.\n", name)
			return
		}

		localPath := localRepo.FindCommandSetFile(name)
		if !rmYes && !confirmPrompt(fmt.Sprintf("Remove local command set '%s'?", name)) {
			fmt.Println("⏭️  Remove cancelled.")
			return
		}

		if err := localRepo.DeleteCommandSet(name); err != nil {
			fmt.Printf("❌ Failed to remove '%s': %v\n", name, err)
			return
		}

		pruneEmptyDirectories(filepath.Dir(localPath), localRepo.GetPath())
		fmt.Printf("✅ Removed local command set '%s'.\n", name)
	},
}

var openCmd = &cobra.Command{
	Use:   "open [command-set-name]",
	Short: "Open a command set file for editing",
	Long: `Open a local command set file in an available editor.
If the command set only exists in bundled repository, ShellDock first copies it to local repository, then opens it.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, err := normalizeCommandSetName(args[0])
		handleError(err)

		manager, err := repo.NewManager()
		handleError(err)

		localRepo := manager.GetLocalRepo()
		localPath := localRepo.FindCommandSetFile(name)

		if localPath == "" {
			bundledRepo := manager.GetBundledRepo()
			if bundledRepo == nil || !bundledRepo.Exists(name) {
				fmt.Printf("❌ Command set '%s' not found.\n", name)
				fmt.Printf("💡 Create one with: shelldock add %s\n", name)
				return
			}

			localPath, err = copyBundledCommandSetToLocal(localRepo, bundledRepo, name)
			if err != nil {
				fmt.Printf("❌ Failed to create local editable copy for '%s': %v\n", name, err)
				return
			}

			fmt.Printf("📥 Copied '%s' to local repository for editing.\n", name)
		}

		fmt.Printf("📝 Opening: %s\n", localPath)
		if err := openInEditor(localPath, openEditor); err != nil {
			fmt.Printf("❌ Failed to open editor: %v\n", err)
		}
	},
}

func init() {
	addCmd.Flags().BoolVar(&addOpen, "open", false, "Open file in editor after creating/copying")
	addCmd.Flags().StringVar(&addDescription, "description", "", "Description for newly created command set")
	addCmd.Flags().StringVar(&addVersion, "version", "v1", "Version for newly created command set")
	addCmd.Flags().StringVarP(&openEditor, "editor", "e", "", "Editor command to use (overrides auto-detection)")

	rmCmd.Flags().BoolVarP(&rmYes, "yes", "y", false, "Skip confirmation prompt")

	openCmd.Flags().StringVarP(&openEditor, "editor", "e", "", "Editor command to use (overrides auto-detection)")
}

func normalizeCommandSetName(input string) (string, error) {
	name := strings.TrimSpace(input)
	name = strings.TrimSuffix(name, ".yaml")
	if name == "" {
		return "", fmt.Errorf("command set name cannot be empty")
	}
	if strings.ContainsAny(name, `/\`) {
		return "", fmt.Errorf("command set name must not include path separators")
	}
	if strings.Contains(name, "..") {
		return "", fmt.Errorf("command set name must not contain '..'")
	}
	return name, nil
}

func openInEditor(filePath, preferred string) error {
	editor, err := selectEditorCommand(preferred)
	if err != nil {
		return err
	}

	args := append(editor.args, filePath)
	cmd := exec.Command(editor.bin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func selectEditorCommand(preferred string) (editorCommand, error) {
	if strings.TrimSpace(preferred) != "" {
		cmd, err := parseEditorCommand(preferred)
		if err != nil {
			return editorCommand{}, err
		}
		if !isCommandAvailable(cmd.bin) {
			return editorCommand{}, fmt.Errorf("requested editor '%s' is not available in PATH", cmd.bin)
		}
		return cmd, nil
	}

	if envEditor := strings.TrimSpace(os.Getenv("EDITOR")); envEditor != "" {
		cmd, err := parseEditorCommand(envEditor)
		if err == nil && isCommandAvailable(cmd.bin) {
			return cmd, nil
		}
	}

	candidates := editorCandidatesByPlatform(runtime.GOOS)
	for _, candidate := range candidates {
		if isCommandAvailable(candidate.bin) {
			return candidate, nil
		}
	}

	return editorCommand{}, fmt.Errorf("no supported editor found for %s (tried: %s)", runtime.GOOS, strings.Join(editorNames(candidates), ", "))
}

func parseEditorCommand(value string) (editorCommand, error) {
	parts := strings.Fields(strings.TrimSpace(value))
	if len(parts) == 0 {
		return editorCommand{}, fmt.Errorf("editor command is empty")
	}
	return editorCommand{
		bin:  parts[0],
		args: parts[1:],
	}, nil
}

func editorCandidatesByPlatform(goos string) []editorCommand {
	common := []editorCommand{
		{bin: "code", args: []string{"--wait"}},
		{bin: "nvim"},
		{bin: "vim"},
		{bin: "nano"},
		{bin: "vi"},
	}

	switch goos {
	case "windows":
		return []editorCommand{
			{bin: "code", args: []string{"--wait"}},
			{bin: "notepad++"},
			{bin: "notepad"},
		}
	case "darwin":
		return append(common, editorCommand{bin: "open", args: []string{"-t"}})
	default:
		return append(common, editorCommand{bin: "micro"})
	}
}

func editorNames(candidates []editorCommand) []string {
	var names []string
	for _, candidate := range candidates {
		names = append(names, candidate.bin)
	}
	return dedupeStrings(names)
}

func dedupeStrings(items []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func pruneEmptyDirectories(startDir, stopDir string) {
	current := startDir
	stopDir = filepath.Clean(stopDir)
	for current != "." {
		rel, err := filepath.Rel(stopDir, current)
		if err != nil || strings.HasPrefix(rel, "..") {
			return
		}
		if current == stopDir {
			return
		}

		entries, err := os.ReadDir(current)
		if err != nil || len(entries) > 0 {
			return
		}

		_ = os.Remove(current)
		current = filepath.Dir(current)
	}
}

func copyBundledCommandSetToLocal(localRepo, bundledRepo *repo.Repository, name string) (string, error) {
	sourcePath := bundledRepo.FindCommandSetFile(name)
	if sourcePath == "" {
		return "", fmt.Errorf("bundled command set '%s' not found", name)
	}

	sourceData, err := os.ReadFile(sourcePath)
	if err != nil {
		return "", fmt.Errorf("failed to read bundled file: %w", err)
	}

	if err := os.MkdirAll(localRepo.GetPath(), 0755); err != nil {
		return "", fmt.Errorf("failed to create local repository: %w", err)
	}

	targetPath := filepath.Join(localRepo.GetPath(), name+".yaml")
	if err := os.WriteFile(targetPath, sourceData, 0644); err != nil {
		return "", fmt.Errorf("failed to write local file: %w", err)
	}

	return targetPath, nil
}
