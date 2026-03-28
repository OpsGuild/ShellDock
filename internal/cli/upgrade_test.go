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

func TestExtractVersionFromText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "shelldock version 1.6", expected: "1.6"},
		{input: "v2.0.1", expected: "v2.0.1"},
		{input: "current=v1.3.0-beta.1", expected: "v1.3.0-beta.1"},
		{input: "no version here", expected: ""},
	}

	for _, tt := range tests {
		got := extractVersionFromText(tt.input)
		if got != tt.expected {
			t.Errorf("extractVersionFromText(%q) = %q, expected %q", tt.input, got, tt.expected)
		}
	}
}

func TestNormalizeVersionForCompare(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "v1.6", expected: "1.6"},
		{input: "1.6", expected: "1.6"},
		{input: " V2.0 ", expected: "2.0"},
	}

	for _, tt := range tests {
		got := normalizeVersionForCompare(tt.input)
		if got != tt.expected {
			t.Errorf("normalizeVersionForCompare(%q) = %q, expected %q", tt.input, got, tt.expected)
		}
	}
}
