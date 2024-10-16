package main

import (
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var testExecutablePath = ""

func ensureExecutable(t *testing.T) string {
	t.Helper()

	if testExecutablePath == "" {
		b, err := exec.Command("go", "build").CombinedOutput()
		if err != nil {
			t.Fatal("Compile error while building `changelog-from-release` executable:", err, ":", string(b))
		}
		if runtime.GOOS == "windows" {
			testExecutablePath = `.\changelog-from-release.exe`
		} else {
			testExecutablePath = `./changelog-from-release`
		}
	}

	return testExecutablePath
}

func TestGenerateChangelog(t *testing.T) {
	exe := ensureExecutable(t)

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
	exe := ensureExecutable(t)

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
	exe := ensureExecutable(t)

	b, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		panic(err)
	}
	want := strings.ReplaceAll(string(b), "\r\n", "\n") // Replace \r\n for Windows

	b, err = exec.Command(exe, "-r", "https://github.com/rhysd/changelog-from-release.git").CombinedOutput()
	have := string(b)
	if err != nil {
		t.Fatal(err, have, exe)
	}

	if want != have {
		t.Fatalf("Generated output was different from CHANGELOG.md\n\n%s", cmp.Diff(have, want))
	}
}

func TestInvalidRemoteURL(t *testing.T) {
	exe := ensureExecutable(t)
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

func TestInvalidGitHubToken(t *testing.T) {
	exe := ensureExecutable(t)

	c := exec.Command(exe)
	c.Env = append(c.Environ(), "GITHUB_TOKEN=invalid")
	b, err := c.CombinedOutput()
	out := string(b)
	if err == nil {
		t.Fatalf("error did not happen: %q", out)
	}
	if !strings.Contains(out, "401 Bad credentials") {
		t.Fatalf("Wanted 401 bad credential error but got %q", out)
	}
}
