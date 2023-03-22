package main

import (
	"fmt"
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

func TestGenerateWithRemoteURL(t *testing.T) {
	exe := validateExecutable(t)
	b, err := exec.Command(exe, "-r", "https://github.com/rhysd/action-setup-vim").CombinedOutput()
	out := string(b)
	if err != nil {
		t.Fatal(err, out, exe)
	}

	for _, v := range []string{
		"v1.2.14",
		"v1.2.0",
		"v1.1.3",
		"v1.1.0",
		"v1.0.2",
		"v1.0.0",
	} {
		var link string
		if v == "v1.0.0" {
			link = "[v1.0.0]: https://github.com/rhysd/action-setup-vim/tree/v1.0.0"
		} else {
			link = fmt.Sprintf("[%s]: https://github.com/rhysd/action-setup-vim/compare/", v)
		}
		for _, want := range []string{
			fmt.Sprintf(`<a name="%s"></a>`, v),
			fmt.Sprintf(`# [%s](https://github.com/rhysd/action-setup-vim/releases/tag/%s)`, v, v),
			fmt.Sprintf(`[Changes][%s]`, v),
			fmt.Sprintf(`[Changes][%s]`, v),
			link,
		} {
			if !strings.Contains(out, want) {
				t.Fatalf("%q was not contained in output:\n%s", want, out)
			}
		}
	}

	r := regexp.MustCompile(`<!-- Generated by https://github.com/rhysd/changelog-from-release v\d+\.\d+\.\d+ -->\n$`)
	if !r.Match(b) {
		t.Fatalf("Footer does not exist at end of input:\n%s", out)
	}
}

func TestInvalidRemoteURL(t *testing.T) {
	exe := validateExecutable(t)
	tests := []struct {
		what  string
		input string
		want  string
	}{
		{
			what:  "invalid URL",
			input: "hello",
			want:  `could not send HEAD request to Git remote URL "hello"`,
		},
		{
			what:  "file URL",
			input: "file:///path/to/file.txt",
			want:  `unsupported protocol scheme "file"`,
		},
		{
			what:  "not a GitHub URL",
			input: "https://example.com",
			want:  "only 'github.com' is supported but got 'https://example.com'",
		},
		{
			what:  "repository does not exist",
			input: "https://github.com/rhysd/this-repository-does-not-exist-oops",
			want:  "404 Not Found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.what, func(t *testing.T) {
			b, err := exec.Command(exe, "-r", tc.input).CombinedOutput()
			out := string(b)
			if err == nil {
				t.Fatal("Error did not occur", out, exe)
			}
			if !strings.Contains(out, tc.want) {
				t.Fatalf("Output %q does not contain %q", out, tc.want)
			}
		})
	}
}
