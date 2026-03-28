package cli

import "testing"

func TestNormalizeCommandSetName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{input: "docker", expected: "docker"},
		{input: " docker ", expected: "docker"},
		{input: "docker.yaml", expected: "docker"},
		{input: "", hasError: true},
		{input: "a/b", hasError: true},
		{input: "..", hasError: true},
	}

	for _, tt := range tests {
		result, err := normalizeCommandSetName(tt.input)
		if tt.hasError {
			if err == nil {
				t.Errorf("normalizeCommandSetName(%q) expected error, got nil", tt.input)
			}
			continue
		}

		if err != nil {
			t.Errorf("normalizeCommandSetName(%q) unexpected error: %v", tt.input, err)
			continue
		}

		if result != tt.expected {
			t.Errorf("normalizeCommandSetName(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestEditorCandidatesByPlatform(t *testing.T) {
	tests := []struct {
		goos         string
		expectedHead string
	}{
		{goos: "linux", expectedHead: "code"},
		{goos: "darwin", expectedHead: "code"},
		{goos: "windows", expectedHead: "code"},
	}

	for _, tt := range tests {
		candidates := editorCandidatesByPlatform(tt.goos)
		if len(candidates) == 0 {
			t.Fatalf("editorCandidatesByPlatform(%q) returned no candidates", tt.goos)
		}
		if candidates[0].bin != tt.expectedHead {
			t.Errorf("editorCandidatesByPlatform(%q) first candidate = %q, expected %q", tt.goos, candidates[0].bin, tt.expectedHead)
		}
	}
}
