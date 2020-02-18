package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSmoke(t *testing.T) {
	for _, tc := range []struct {
		args   []string
		stdout string
	}{
		{
			[]string{},
			"",
		},
		{
			[]string{"-t"},
			version + "\n",
		},
		{
			[]string{"-v"},
			version + "\n",
		},
	} {
		t.Run(fmt.Sprint(tc.args), func(t *testing.T) {
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

			p := filepath.Join(cwd, exe)
			b, err := exec.Command(p, tc.args...).CombinedOutput()
			out := string(b)
			if err != nil {
				t.Fatal(err, out, p)
			}
			if out != tc.stdout {
				t.Fatalf("Should output %#v but got %#v", tc.stdout, out)
			}

			git, err := NewGitForCwd()
			if err != nil {
				t.Fatal(err)
			}

			if out, err = git.Exec("diff", "--quiet"); err != nil {
				t.Fatal("Repository should not be dirty:", out, err)
			}
		})
	}
}
