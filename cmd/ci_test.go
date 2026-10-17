package cmd

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestCIWorkflowHardening(t *testing.T) {
	path := filepath.Join("..", ".github", "workflows", "ci.yml")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read %s: %v", path, err)
	}
	content := string(raw)

	// 1. Top-level permissions: {}
	if !strings.Contains(content, "permissions: {}") {
		t.Error("ci.yml must specify top-level permissions: {} for least-privilege token access")
	}

	// 2. Concurrency block
	if !strings.Contains(content, "concurrency:") || !strings.Contains(content, "cancel-in-progress: true") {
		t.Error("ci.yml must define concurrency group with cancel-in-progress: true")
	}

	// 3. Timeout on jobs
	if !strings.Contains(content, "timeout-minutes:") {
		t.Error("ci.yml jobs must declare timeout-minutes")
	}

	// 4. Action SHA-pinning (no mutable @vX tags for actions/*)
	actionRegex := regexp.MustCompile(`uses:\s*(actions/[a-zA-Z0-9_-]+)@([^\s#]+)`)
	matches := actionRegex.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		t.Fatal("no action references found in ci.yml")
	}
	sha40 := regexp.MustCompile(`^[0-9a-f]{40}$`)
	for _, match := range matches {
		actionName := match[1]
		ref := match[2]
		if !sha40.MatchString(ref) {
			t.Errorf("action %s pinned to %q; must be pinned to 40-character commit SHA", actionName, ref)
		}
	}

	// 5. Checkout persist-credentials: false
	if !strings.Contains(content, "persist-credentials: false") {
		t.Error("actions/checkout step must set persist-credentials: false")
	}

	// 6. gofmt formatting check
	if !strings.Contains(content, "gofmt") {
		t.Error("ci.yml must include gofmt check")
	}

	// 7. govulncheck dependency vulnerability audit
	if !strings.Contains(content, "govulncheck") {
		t.Error("ci.yml must include govulncheck vulnerability audit")
	}

	// 8. shellcheck for scripts
	if !strings.Contains(content, "shellcheck") {
		t.Error("ci.yml must include shellcheck for shell scripts")
	}
}
