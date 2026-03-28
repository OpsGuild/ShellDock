package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type upgradeMethod string

const (
	upgradeMethodAPT          upgradeMethod = "apt"
	upgradeMethodDNF          upgradeMethod = "dnf"
	upgradeMethodYUM          upgradeMethod = "yum"
	upgradeMethodPacman       upgradeMethod = "pacman"
	upgradeMethodYay          upgradeMethod = "yay"
	upgradeMethodParu         upgradeMethod = "paru"
	upgradeMethodBrew         upgradeMethod = "brew"
	upgradeMethodSnap         upgradeMethod = "snap"
	upgradeMethodFlatpak      upgradeMethod = "flatpak"
	upgradeMethodChocolatey   upgradeMethod = "chocolatey"
	upgradeMethodDirectBinary upgradeMethod = "direct-binary"
)

type upgradePlan struct {
	method   upgradeMethod
	reason   string
	commands []string
}

type githubRelease struct {
	TagName string `json:"tag_name"`
}

var (
	upgradeYes    bool
	upgradeDryRun bool
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the ShellDock binary to latest version",
	Long: `Upgrade ShellDock itself to the latest available version.
The command detects your platform and likely installation method, then runs the appropriate upgrade workflow.`,
	Run: func(cmd *cobra.Command, args []string) {
		plan := detectUpgradePlan()

		fmt.Printf("🔍 Detected install method: %s\n", plan.method)
		fmt.Printf("🖥️  Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		if strings.TrimSpace(plan.reason) != "" {
			fmt.Printf("ℹ️  %s\n", plan.reason)
		}

		if plan.method == upgradeMethodDirectBinary {
			assetName, err := releaseAssetNameForPlatform(runtime.GOOS, runtime.GOARCH)
			if err != nil {
				fmt.Printf("❌ %v\n", err)
				return
			}
			fmt.Printf("📦 Upgrade asset: %s\n", assetName)
		} else {
			fmt.Println("📋 Upgrade commands:")
			for _, command := range plan.commands {
				fmt.Printf("   $ %s\n", command)
			}
		}

		if upgradeDryRun {
			fmt.Println("🧪 Dry run mode enabled. No changes made.")
			return
		}

		if !upgradeYes && !confirmPrompt("Proceed with ShellDock upgrade?") {
			fmt.Println("⏭️  Upgrade cancelled.")
			return
		}

		var err error
		if plan.method == upgradeMethodDirectBinary {
			err = runDirectBinaryUpgrade()
		} else {
			err = runUpgradeCommands(plan.commands)
		}
		if err != nil {
			fmt.Printf("❌ Upgrade failed: %v\n", err)
			return
		}

		fmt.Println("✅ ShellDock upgrade complete.")
		fmt.Println("💡 Run 'shelldock --version' to verify the installed version.")
	},
}

func init() {
	upgradeCmd.Flags().BoolVarP(&upgradeYes, "yes", "y", false, "Skip confirmation prompt")
	upgradeCmd.Flags().BoolVar(&upgradeDryRun, "dry-run", false, "Show detected method and commands without executing")
}

func detectUpgradePlan() upgradePlan {
	exePath, _ := os.Executable()
	exePath = filepath.ToSlash(exePath)
	goos := runtime.GOOS

	if goos == "windows" {
		if isChocolateyInstalled() {
			return upgradePlan{
				method:   upgradeMethodChocolatey,
				reason:   "Detected ShellDock installed via Chocolatey.",
				commands: []string{"choco upgrade shelldock -y"},
			}
		}
		return upgradePlan{
			method: upgradeMethodDirectBinary,
			reason: "Chocolatey installation not detected. Falling back to direct binary upgrade.",
		}
	}

	if strings.Contains(exePath, "/snap/") || commandSucceeds("snap", "list", "shelldock") {
		return upgradePlan{
			method:   upgradeMethodSnap,
			reason:   "Detected snap-managed installation.",
			commands: []string{"sudo snap refresh shelldock"},
		}
	}

	if strings.Contains(exePath, "/flatpak/") || commandSucceeds("flatpak", "info", "com.github.opsguild.shelldock") {
		return upgradePlan{
			method:   upgradeMethodFlatpak,
			reason:   "Detected flatpak-managed installation.",
			commands: []string{"flatpak update -y com.github.opsguild.shelldock"},
		}
	}

	if isHomebrewInstalled(exePath) {
		return upgradePlan{
			method:   upgradeMethodBrew,
			reason:   "Detected Homebrew-managed installation.",
			commands: []string{"brew update", "brew upgrade shelldock"},
		}
	}

	if isAPTInstalled() {
		return upgradePlan{
			method:   upgradeMethodAPT,
			reason:   "Detected APT-managed installation.",
			commands: []string{"sudo apt update", "sudo apt install --only-upgrade -y shelldock"},
		}
	}

	if isRPMInstalled() {
		if isCommandAvailable("dnf") {
			return upgradePlan{
				method:   upgradeMethodDNF,
				reason:   "Detected RPM package with DNF available.",
				commands: []string{"sudo dnf upgrade -y shelldock"},
			}
		}
		if isCommandAvailable("yum") {
			return upgradePlan{
				method:   upgradeMethodYUM,
				reason:   "Detected RPM package with YUM available.",
				commands: []string{"sudo yum update -y shelldock"},
			}
		}
	}

	if isAURInstalledWith("yay") {
		return upgradePlan{
			method:   upgradeMethodYay,
			reason:   "Detected AUR package installed with yay.",
			commands: []string{"yay -Syu --noconfirm shelldock"},
		}
	}

	if isAURInstalledWith("paru") {
		return upgradePlan{
			method:   upgradeMethodParu,
			reason:   "Detected AUR package installed with paru.",
			commands: []string{"paru -Syu --noconfirm shelldock"},
		}
	}

	if isPacmanInstalled() {
		return upgradePlan{
			method:   upgradeMethodPacman,
			reason:   "Detected pacman-managed installation.",
			commands: []string{"sudo pacman -Syu --noconfirm shelldock"},
		}
	}

	return upgradePlan{
		method: upgradeMethodDirectBinary,
		reason: "No package-manager ownership detected. Falling back to direct binary upgrade.",
	}
}

func runUpgradeCommands(commands []string) error {
	for _, command := range commands {
		fmt.Printf("➡️  Running: %s\n", command)
		if err := runShellCommand(command); err != nil {
			return err
		}
	}
	return nil
}

func runDirectBinaryUpgrade() error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("direct-binary auto-upgrade on Windows is not supported; use 'choco upgrade shelldock' or manual release download")
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to detect executable path: %w", err)
	}

	tag, err := fetchLatestReleaseTag()
	if err != nil {
		return fmt.Errorf("failed to detect latest release: %w", err)
	}

	assetName, err := releaseAssetNameForPlatform(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return err
	}

	downloadURL := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", githubRepo, tag, assetName)
	fmt.Printf("📥 Downloading %s...\n", downloadURL)

	downloadedFile, err := downloadReleaseBinary(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download release binary: %w", err)
	}
	defer func() { _ = os.Remove(downloadedFile) }()

	destDir := filepath.Dir(exePath)
	stagedPath := filepath.Join(destDir, ".shelldock.new")
	_ = os.Remove(stagedPath)

	if err := copyFileWithMode(downloadedFile, stagedPath, 0755); err == nil {
		if err := os.Rename(stagedPath, exePath); err == nil {
			fmt.Printf("✅ Upgraded ShellDock to %s at %s\n", tag, exePath)
			return nil
		}
		_ = os.Remove(stagedPath)
	}

	if !isCommandAvailable("sudo") {
		return fmt.Errorf("insufficient permissions to replace binary at %s; re-run with elevated privileges", exePath)
	}

	installCmd := fmt.Sprintf("sudo install -m 755 %s %s", shellQuote(downloadedFile), shellQuote(exePath))
	fmt.Println("🔐 Escalating with sudo to replace the binary...")
	if err := runShellCommand(installCmd); err != nil {
		return fmt.Errorf("failed to replace binary with sudo: %w", err)
	}

	fmt.Printf("✅ Upgraded ShellDock to %s at %s\n", tag, exePath)
	return nil
}

func releaseAssetNameForPlatform(goos, goarch string) (string, error) {
	archMap := map[string]string{
		"x86_64":  "amd64",
		"amd64":   "amd64",
		"aarch64": "arm64",
		"arm64":   "arm64",
	}

	arch, ok := archMap[goarch]
	if !ok {
		return "", fmt.Errorf("unsupported architecture for direct upgrade: %s", goarch)
	}

	switch goos {
	case "linux", "darwin":
		return fmt.Sprintf("shelldock-%s-%s", goos, arch), nil
	case "windows":
		return fmt.Sprintf("shelldock-windows-%s.exe", arch), nil
	default:
		return "", fmt.Errorf("unsupported operating system for direct upgrade: %s", goos)
	}
}

func fetchLatestReleaseTag() (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodGet, githubAPI+"/releases/latest", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "shelldock-upgrade")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github API returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var release githubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return "", err
	}
	if strings.TrimSpace(release.TagName) == "" {
		return "", fmt.Errorf("latest release has empty tag name")
	}

	return release.TagName, nil
}

func downloadReleaseBinary(url string) (string, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	file, err := os.CreateTemp("", "shelldock-upgrade-*")
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(file, resp.Body); err != nil {
		_ = file.Close()
		_ = os.Remove(file.Name())
		return "", err
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(file.Name())
		return "", err
	}
	if err := os.Chmod(file.Name(), 0755); err != nil {
		_ = os.Remove(file.Name())
		return "", err
	}

	return file.Name(), nil
}

func copyFileWithMode(src, dst string, mode os.FileMode) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = source.Close() }()

	destination, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer func() { _ = destination.Close() }()

	if _, err := io.Copy(destination, source); err != nil {
		return err
	}

	return destination.Sync()
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `'\''`) + "'"
}

func runShellCommand(command string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-Command", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func commandSucceeds(name string, args ...string) bool {
	if !isCommandAvailable(name) {
		return false
	}
	cmd := exec.Command(name, args...)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd.Run() == nil
}

func isChocolateyInstalled() bool {
	if !isCommandAvailable("choco") {
		return false
	}
	cmd := exec.Command("choco", "list", "--local-only", "--exact", "shelldock", "--limit-output")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	output := strings.ToLower(string(out))
	return strings.Contains(output, "shelldock|")
}

func isHomebrewInstalled(exePath string) bool {
	if strings.Contains(exePath, "/cellar/shelldock/") || strings.Contains(exePath, "/homebrew/") {
		return true
	}
	if commandSucceeds("brew", "list", "--formula", "shelldock") {
		return true
	}
	return commandSucceeds("brew", "list", "--formula", "opsguild/tap/shelldock")
}

func isAPTInstalled() bool {
	return commandSucceeds("dpkg", "-s", "shelldock")
}

func isRPMInstalled() bool {
	return commandSucceeds("rpm", "-q", "shelldock")
}

func isPacmanInstalled() bool {
	return commandSucceeds("pacman", "-Q", "shelldock")
}

func isAURInstalledWith(tool string) bool {
	return commandSucceeds(tool, "-Q", "shelldock")
}
