package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

// ResolveRedirect resolves URL redirects and returns parsed URL
func ResolveRedirect(u string) (*url.URL, error) {
	u = strings.TrimSuffix(u, ".git")

	res, err := http.Head(u)
	if err != nil {
		return nil, fmt.Errorf("could not send HEAD request to Git remote URL %q for following repository redirect: %w", u, err)
	}

	// GitHub returns 404 when the repository is private. GHE would do the same since all GHE repositories
	// basically require authentication. 403 may be returned as well. (#19)
	if res.StatusCode != 200 && res.StatusCode != 404 && res.StatusCode != 403 {
		return nil, fmt.Errorf("HEAD request to Git remote URL %q for following repository redirect was not successful: %s", u, res.Status)
	}

	// GitHub Enterprise server redirects the repository access to the login page URL with 200 status.
	if res.Request.URL.Path == "/login" && res.Request.URL.RawQuery != "" {
		slog.Debug("URL was redirected to login page. Gave up resolving the redirect", "from", u, "to", res.Request.URL)
		return url.Parse(u)
	}

	slog.Debug("Resolved URL", "from", u, "to", res.Request.URL)
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
	slog.Debug("Running Git command", "bin", git.bin, "args", a)
	return exec.Command(git.bin, a...)
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

	s := string(out)
	slog.Debug("Git command successfully exited", "output", s)
	return s, nil
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

	slog.Debug("Extracted the first remote name", "name", s)
	return s, nil
}

// FirstRemoteURL returns a URL of remote repository. When multiple remotes are configured, the
// first one will be chosen.
func (git *Git) FirstRemoteURL() (*url.URL, error) {
	r, err := git.FirstRemoteName()
	if err != nil {
		return nil, fmt.Errorf("could not get URL of remote repository: %w", err)
	}

	c := fmt.Sprintf("remote.%s.url", r)
	s, err := git.Exec("config", c)
	if err != nil {
		return nil, fmt.Errorf("could not get URL of remote %q: %w", r, err)
	}
	slog.Debug("Got config for remote URL", "config", c, "url", s)

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
	slog.Debug("Converted to HTTP URL", "url", s)

	u, err := ResolveRedirect(s)
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
	slog.Debug("Git execution environment", "executable", exe, "cwd", cwd)
	return &Git{exe, cwd}, nil
}
