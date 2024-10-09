package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/go-github/v65/github"
	"golang.org/x/oauth2"
)

// GitHub implements GitHub API v3 client
type GitHub struct {
	api      *github.Client
	apiCtx   context.Context
	owner    string
	repoName string
}

// Releases fetches releases information. When no release is found, this method returns an error
func (gh *GitHub) Releases() ([]*github.RepositoryRelease, error) {
	rels := []*github.RepositoryRelease{}
	page := 1
	for {
		opts := github.ListOptions{
			Page:    page,
			PerPage: 100,
		}
		rs, res, err := gh.api.Repositories.ListReleases(gh.apiCtx, gh.owner, gh.repoName, &opts)
		if err != nil {
			return nil, fmt.Errorf("cannot get releases from repository %s/%s via GitHub API: %w", gh.owner, gh.repoName, err)
		}
		rels = append(rels, rs...)
		if res.NextPage == 0 {
			return rels, nil
		}
		page = res.NextPage
	}
}

func (gh *GitHub) CustomAutolinks() ([]*github.Autolink, error) {
	links := []*github.Autolink{}
	page := 1
	for {
		opts := github.ListOptions{Page: page}
		ls, res, err := gh.api.Repositories.ListAutolinks(gh.apiCtx, gh.owner, gh.repoName, &opts)
		if err != nil {
			return nil, fmt.Errorf("cannot get custom autolinks from repository %s/%s via GitHub API: %w", gh.owner, gh.repoName, err)
		}
		links = append(links, ls...)
		if res.NextPage == 0 {
			return links, nil
		}
		page = res.NextPage
	}
}

// NewGitHub creates GitHub instance from given repository URL
func NewGitHub(u *url.URL, c context.Context) (*GitHub, error) {
	// '/owner/name'
	path := strings.TrimSuffix(u.Path, ".git")
	slug := strings.Split(path, "/")
	if len(slug) != 3 {
		return nil, fmt.Errorf("invalid slug in GitHub repository URL path: %s", path)
	}

	client := http.DefaultClient
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		client = oauth2.NewClient(c, src)
	}

	api := github.NewClient(client)
	if v := os.Getenv("GITHUB_API_BASE_URL"); v != "" {
		// > BaseURL should always be specified with a trailing slash.
		// https://pkg.go.dev/github.com/google/go-github/github#Client
		if !strings.HasSuffix(v, "/") {
			v += "/"
		}

		u, err := url.Parse(v)
		if err != nil {
			return nil, fmt.Errorf("invalid URL %q in $GITHUB_API_BASE_URL: %w", v, err)
		}

		api.BaseURL = u
	}
	return &GitHub{api, c, slug[1], slug[2]}, nil
}
