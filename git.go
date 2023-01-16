package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func resolveRedirect(url string) (*url.URL, error) {
	url = strings.TrimSuffix(url, ".git")

	res, err := http.Head(url)
	if err != nil {
		return nil, fmt.Errorf("could not send HEAD request to Git remote URL %s for following repository redirect: %w", url, err)
	}

	// Do not check response status because it is not useful. GitHub returns 404 when the repository is private.
	// And GHE would does the same since all GHE repositories basically require authentication. (#19)
	return res.Request.URL, nil
}

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
	out, err := git.Command(subcmd, args...).CombinedOutput()
	out = bytes.TrimSpace(out)

	if err != nil {
		for i, c := range out {
			if c == '\n' || c == '\r' {
				out[i] = ' '
			}
		}
		return "", fmt.Errorf("Git command %q  with args %v failed with output %q: %w", subcmd, args, out, err)
	}

	return string(out), nil
}

// FirstRemoteName returns remote name of current Git repository. When multiple remotes are
// configured, the first one will be chosen.
func (git *Git) FirstRemoteName() (string, error) {
	s, err := git.Exec("remote")
	if err != nil {
		return "", fmt.Errorf("could not retrieve remote name: %w", err)
	}

	if i := strings.IndexAny(s, "\r\n"); i >= 0 {
		s = s[:i]
	}

	if s == "" {
		return "", fmt.Errorf("no remote is configured in this repository")
	}

	return s, nil
}

// FirstRemoteURL returns a URL of remote repository. When multiple remotes are configured, the
// first one will be chosen.
func (git *Git) FirstRemoteURL() (*url.URL, error) {
	r, err := git.FirstRemoteName()
	if err != nil {
		return nil, fmt.Errorf("could not get URL of remote repository: %w", err)
	}

	s, err := git.Exec("config", fmt.Sprintf("remote.%s.url", r))
	if err != nil {
		return nil, fmt.Errorf("could not get URL of remote %q: %w", r, err)
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

	u, err := resolveRedirect(s)
	if err != nil {
		return nil, err
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
