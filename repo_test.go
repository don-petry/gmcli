package main

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// Third-party actions must stay pinned to full commit SHAs so a moved tag
// cannot change what CI runs.
func TestWorkflowActionsPinnedToSHA(t *testing.T) {
	files, err := filepath.Glob(filepath.Join(".github", "workflows", "*.yml"))
	if err != nil || len(files) == 0 {
		t.Fatalf("no workflow files found: %v", err)
	}
	uses := regexp.MustCompile(`(?m)^\s*-?\s*uses:\s*([^\s#]+)`)
	pinned := regexp.MustCompile(`@[0-9a-f]{40}$`)
	for _, file := range files {
		raw, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range uses.FindAllStringSubmatch(string(raw), -1) {
			if ref := m[1]; ref[0] != '.' && !pinned.MatchString(ref) {
				t.Errorf("%s: %s is not pinned to a commit SHA", file, ref)
			}
		}
	}
}
