package main

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type GitHub struct {
	api      *github.Client
	apiCtx   context.Context
	owner    string
	repoName string
}

func (gh *GitHub) Releases() ([]*github.RepositoryRelease, error) {
	rels := []*github.RepositoryRelease{}
	for {
		rs, res, err := gh.api.Repositories.ListReleases(gh.apiCtx, gh.owner, gh.repoName, nil)
		if err != nil {
			return nil, errors.Wrapf(err, "Cannot get releases from repository %s/%s via GitHub API", gh.owner, gh.repoName)
		}
		rels = append(rels, rs...)
		if res.NextPage == 0 {
			return rels, nil
		}
	}
}

func GitHubFromURL(u *url.URL) (*GitHub, error) {
	if u.Host != "github.com" {
		return nil, errors.Errorf("Only 'github.com' is supported but got '%s'", u.String())
	}

	// '/owner/name'
	path := strings.TrimSuffix(u.Path, ".git")
	slug := strings.Split(path, "/")
	if len(slug) != 3 {
		return nil, errors.Errorf("Invalid slug of GitHub repo: %s", path)
	}

	ctx := context.Background()
	client := http.DefaultClient
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		client = oauth2.NewClient(ctx, src)
	}
	return &GitHub{github.NewClient(client), ctx, slug[1], slug[2]}, nil
}
