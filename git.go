package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

// Git represents Git command for specific repository
type Git struct {
	bin  string
	root string
}

// Command returns exec.Command instance which runs given Git subcommand with given arguments
func (git *Git) Command(subcmd string, args ...string) *exec.Cmd {
	// e.g. 'git diff --cached' -> 'git -C /path/to/repo diff --cached'
	a := append([]string{"-C", git.root, subcmd}, args...)
	cmd := exec.Command(git.bin, a...)
	return cmd
}

// Exec runs runs given Git subcommand with given arguments
func (git *Git) Exec(subcmd string, args ...string) (string, error) {
	b, err := git.Command(subcmd, args...).CombinedOutput()

	// Chop last newline
	l := len(b)
	if l > 0 && b[l-1] == '\n' {
		b = b[:l-1]
	}

	// Make output in oneline in error cases
	for i, c := range b {
		if c == '\n' {
			b[i] = ' '
		}
	}
	out := string(b)

	if err != nil {
		return "", fmt.Errorf("Git command %q %v failed with output %q: %w", subcmd, args, out, err)
	}

	return out, nil
}

// FirstRemoteName returns remote name of current Git repository. When multiple remotes are
// configured, the first one will be chosen. When no remote is configured, this method returns empty
// string (no error is returned).
func (git *Git) FirstRemoteName() (string, error) {
	s, err := git.Exec("remote")
	if err != nil {
		return "", fmt.Errorf("could not retrieve remote name: %s: %w", s, err)
	}

	s = strings.TrimSpace(s)
	if i := strings.IndexAny(s, "\r\n"); i >= 0 {
		s = s[:i]
	}

	return s, nil
}

// RemoteURL returns a URL which the current repository is tracking as remote
func (git *Git) RemoteURL() (*url.URL, error) {
	r, err := git.FirstRemoteName()
	if err != nil {
		return nil, err
	}
	if r == "" {
		return nil, fmt.Errorf("could not get any remote URL since no remote is configured in this repository")
	}

	s, err := git.Exec("config", fmt.Sprintf("remote.%s.url", r))
	if err != nil {
		return nil, fmt.Errorf("could not get URL of remote '%s': %s: %w", r, s, err)
	}

	if strings.HasPrefix(s, "git@") && strings.ContainsRune(s, ':') {
		// git@github.com:user/repo.git → https://github.com/user/repo.git
		s = "https://" + strings.Replace(strings.TrimPrefix(s, "git@"), ":", "/", 1)
	} else if strings.HasPrefix(s, "ssh://git@") {
		// ssh://git@github.com/user/repo.git → https://github.com/user/repo.git
		s = "https://" + strings.TrimPrefix(s, "ssh://git@")
	}

	if !strings.HasPrefix(s, "https://") && !strings.HasPrefix(s, "http://") {
		return nil, fmt.Errorf("repository URL is neither HTTP nor HTTPS: %s", s)
	}

	u, err := url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("cannot parse tracking remote URL: %s: %w", s, err)
	}
	return u, nil
}

// NewGitForCwd creates Git instance from Config value. Home directory is assumed to be a root of Git repository
func NewGitForCwd() (*Git, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("cannot get cwd: %w", err)
	}
	exe, err := exec.LookPath("git")
	if err != nil {
		return nil, fmt.Errorf("'git' executable not found: %w", err)
	}
	return &Git{exe, cwd}, nil
}
