package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestSmoke(t *testing.T) {
	exe := "changelog-from-release"
	if _, err := os.Stat(exe); err != nil {
		t.Fatal("Executable not found:", exe)
	}

	if s, err := os.Stat(".git"); err != nil || !s.IsDir() {
		t.Fatal("Test did not run at root of repository")
	}

	b, err := exec.Command(exe).CombinedOutput()
	out := string(b)
	if err != nil {
		t.Fatal(err, out)
	}
	if out != "" {
		t.Fatalf("Should output nothing %#v", out)
	}

	git, err := NewGitForCwd()
	if err != nil {
		t.Fatal(err)
	}

	if out, err = git.Exec("diff", "--quiet"); err != nil {
		t.Fatal("Repository should not be dirty:", out, err)
	}
}
