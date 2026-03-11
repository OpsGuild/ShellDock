package repo

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRepositoryCommandsAvoidEchoDashE(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve test file path")
	}

	repositoryDir := filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..", "repository"))

	err := filepath.WalkDir(repositoryDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || filepath.Ext(path) != ".yaml" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if strings.Contains(string(data), "echo -e") {
			t.Errorf("%s contains echo -e, which is not portable under sh -c", path)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("failed to scan repository command files: %v", err)
	}
}
