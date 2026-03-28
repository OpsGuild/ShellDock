package cli

import "testing"

func TestReleaseAssetNameForPlatform(t *testing.T) {
	tests := []struct {
		goos     string
		goarch   string
		expected string
		hasError bool
	}{
		{goos: "linux", goarch: "amd64", expected: "shelldock-linux-amd64"},
		{goos: "darwin", goarch: "arm64", expected: "shelldock-darwin-arm64"},
		{goos: "windows", goarch: "amd64", expected: "shelldock-windows-amd64.exe"},
		{goos: "linux", goarch: "386", hasError: true},
		{goos: "freebsd", goarch: "amd64", hasError: true},
	}

	for _, tt := range tests {
		result, err := releaseAssetNameForPlatform(tt.goos, tt.goarch)
		if tt.hasError {
			if err == nil {
				t.Errorf("releaseAssetNameForPlatform(%q, %q) expected error, got nil", tt.goos, tt.goarch)
			}
			continue
		}

		if err != nil {
			t.Errorf("releaseAssetNameForPlatform(%q, %q) unexpected error: %v", tt.goos, tt.goarch, err)
			continue
		}

		if result != tt.expected {
			t.Errorf("releaseAssetNameForPlatform(%q, %q) = %q, expected %q", tt.goos, tt.goarch, result, tt.expected)
		}
	}
}

func TestShellQuote(t *testing.T) {
	input := "path/with 'quote'"
	expected := "'path/with '\\''quote'\\'''"
	got := shellQuote(input)
	if got != expected {
		t.Errorf("shellQuote(%q) = %q, expected %q", input, got, expected)
	}
}
