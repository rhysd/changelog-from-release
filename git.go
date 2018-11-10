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
	if err != nil {
		for i := range b {
			if b[i] == '\n' {
				b[i] = ' '
			}
		}
	}

	return string(b), err
}

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
		s = "https://github.com/" + s[len("git@github.com:"):]
	}

	u, err := url.Parse(s)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot parse tracking remote URL: %s", s)
	}
	return u, nil
}

// NewGit creates Git instance from Config value. Home directory is assumed to be a root of Git repository
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
