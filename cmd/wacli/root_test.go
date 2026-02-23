package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestExecute_Version(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	// Execute "version" command
	// We pass stdout/stderr, hoping the command writes to them.
	err := execute([]string{"version"}, stdout, stderr)
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, version) {
		t.Errorf("expected output to contain version %q, got %q", version, out)
	}
}
