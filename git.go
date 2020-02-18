package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
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
		return "", errors.Wrapf(err, "Git command %q %v failed with output %q", subcmd, args, out)
	}

	return out, nil
}

// TrackingRemoteURL returns a URL which the current repository is tracking as remote
func (git *Git) TrackingRemoteURL() (*url.URL, error) {
	s, err := git.Exec("rev-parse", "--abbrev-ref", "--symbolic", "@{u}")
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot retrieve remote name: %s", s)
	}

	// e.g. origin/master
	ss := strings.Split(s, "/")

	if s, err = git.Exec("config", fmt.Sprintf("remote.%s.url", ss[0])); err != nil {
		return nil, errors.Wrapf(err, "Could not get URL of remote '%s': %s", ss[0], s)
	}

	if strings.HasPrefix(s, "git@github.com:") {
		s = "https://github.com/" + strings.TrimPrefix(s, "git@github.com:")
	} else if strings.HasPrefix(s, "ssh://git@github.com/") {
		s = "https://github.com/" + strings.TrimPrefix(s, "ssh://git@github.com/")
	}

	u, err := url.Parse(s)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot parse tracking remote URL: %s", s)
	}
	return u, nil
}

// CheckClean checks if the working tree and index are not dirty
func (git *Git) CheckClean() error {
	if _, err := git.Exec("diff", "--quiet"); err != nil {
		return errors.New("Git working tree is dirty. Please ensure all changes are committed")
	}
	if _, err := git.Exec("diff", "--cached", "--quiet"); err != nil {
		return errors.New("Git index is dirty. Please ensure all changes are added and committed")
	}
	return nil
}

// Add adds given file to Git index tree
func (git *Git) Add(file string) error {
	_, err := git.Exec("add", file)
	return err
}

// Commit creates a new Git commit with given message
func (git *Git) Commit(msg string) error {
	_, err := git.Exec("commit", "-m", msg)
	return err
}

// NewGitForCwd creates Git instance from Config value. Home directory is assumed to be a root of Git repository
func NewGitForCwd() (*Git, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get cwd")
	}
	exe, err := exec.LookPath("git")
	if err != nil {
		return nil, errors.Wrap(err, "'git' executable not found")
	}
	return &Git{exe, cwd}, nil
}
