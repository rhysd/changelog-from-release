package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
)

func TestSmoke(t *testing.T) {
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

	for _, tc := range []struct {
		args   []string
		stdout string
	}{
		{
			[]string{},
			"",
		},
		{
			[]string{"-version"},
			`v\d+\.\d+\.\d+\n$`,
		},
	} {
		t.Run(fmt.Sprint(tc.args), func(t *testing.T) {
			p := filepath.Join(cwd, exe)
			b, err := exec.Command(p, tc.args...).CombinedOutput()
			out := string(b)
			if err != nil {
				t.Fatal(err, out, p)
			}

			re := regexp.MustCompile(tc.stdout)
			if !re.MatchString(out) {
				t.Fatalf("Output %#v did not match to %#v", out, tc.stdout)
			}

			git, err := NewGitForCwd()
			if err != nil {
				t.Fatal(err)
			}

			if err = git.CheckClean(); err != nil {
				t.Fatal("Repository should not be dirty:", err)
			}
		})
	}
}
