package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func validateExecutable(t *testing.T) string {
	exe := "changelog-from-release"
	if runtime.GOOS == "windows" {
		exe = exe + ".exe"
	}

	if _, err := os.Stat(exe); err != nil {
		t.Fatal("Executable not found:", exe)
	}

	if s, err := os.Stat(".git"); err != nil || !s.IsDir() {
		t.Fatal("Test did not run at root of repository")
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(cwd, exe)
}

func TestGenerateChangelog(t *testing.T) {
	exe := validateExecutable(t)

	b, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		panic(err)
	}
	want := strings.ReplaceAll(string(b), "\r\n", "\n") // Replace \r\n for Windows

	b, err = exec.Command(exe).CombinedOutput()
	have := string(b)
	if err != nil {
		t.Fatal(err, have, exe)
	}

	if want != have {
		t.Fatalf("Generated output was different from CHANGELOG.md\n\n%s", cmp.Diff(have, want))
	}
}

func TestVersion(t *testing.T) {
	exe := validateExecutable(t)

	b, err := exec.Command(exe, "-v").CombinedOutput()
	out := string(b)
	if err != nil {
		t.Fatal(err, out, exe)
	}

	r := regexp.MustCompile(`^v\d+\.\d+\.\d+\n$`)
	if !r.Match(b) {
		t.Fatalf("Version output is unexpected: %q", out)
	}
}
