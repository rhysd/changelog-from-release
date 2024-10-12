package main

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v66/github"
)

func TestConfigFilterReleases(t *testing.T) {
	release := func(draft, prerelease bool, tag, desc string) *github.RepositoryRelease {
		return &github.RepositoryRelease{
			Draft:      &draft,
			Prerelease: &prerelease,
			Name:       &desc,
			TagName:    &tag,
		}
	}

	releases := []*github.RepositoryRelease{
		release(false, false, "no-draft-no-prerel", "no draft, no prerelease"),
		release(true, false, "draft-no-prerel", "draft, no prerelease"),
		release(false, true, "no-draft-prerel", "no draft, prerelease"),
		release(true, true, "draft-prerel", "draft, prerelease"),
	}

	tests := []struct {
		cfg  Config
		want []string
	}{
		{
			cfg: Config{},
			want: []string{
				"no draft, no prerelease",
			},
		},
		{
			cfg: Config{
				Drafts: true,
			},
			want: []string{
				"no draft, no prerelease",
				"draft, no prerelease",
			},
		},
		{
			cfg: Config{
				Prerelease: true,
			},
			want: []string{
				"no draft, no prerelease",
				"no draft, prerelease",
			},
		},
		{
			cfg: Config{
				Drafts:     true,
				Prerelease: true,
			},
			want: []string{
				"no draft, no prerelease",
				"draft, no prerelease",
				"no draft, prerelease",
				"draft, prerelease",
			},
		},
		{
			cfg: Config{
				Drafts:     true,
				Prerelease: true,
				Ignore:     regexp.MustCompile(`^draft-`),
			},
			want: []string{
				"no draft, no prerelease",
				"no draft, prerelease",
			},
		},
		{
			cfg: Config{
				Drafts:     true,
				Prerelease: true,
				Extract:    regexp.MustCompile(`^draft-`),
			},
			want: []string{
				"draft, no prerelease",
				"draft, prerelease",
			},
		},
		{
			cfg: Config{
				Drafts:     true,
				Prerelease: true,
				Ignore:     regexp.MustCompile(`.*`),
			},
			want: nil,
		},
		{
			cfg: Config{
				Drafts:     true,
				Prerelease: true,
				Extract:    regexp.MustCompile(`this-regex-never-matches`),
			},
			want: nil,
		},
		{
			cfg: Config{
				Drafts:     true,
				Prerelease: true,
				Ignore:     regexp.MustCompile(`^draft-`),
				Extract:    regexp.MustCompile(`-no-prerel$`),
			},
			want: []string{
				"no draft, no prerelease",
			},
		},
		{
			cfg: Config{
				Drafts:     true,
				Prerelease: true,
				Ignore:     regexp.MustCompile(`this-does-not-match`),
				Extract:    regexp.MustCompile(`.*`),
			},
			want: []string{
				"no draft, no prerelease",
				"draft, no prerelease",
				"no draft, prerelease",
				"draft, prerelease",
			},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%#v", tc.cfg), func(t *testing.T) {
			rels := append([]*github.RepositoryRelease{}, releases...) // t.Run runs tests in parallel
			var have []string
			for _, r := range tc.cfg.filterReleases(rels) {
				have = append(have, r.GetName())
			}

			if !cmp.Equal(have, tc.want) {
				t.Fatal(cmp.Diff(have, tc.want), have)
			}
		})
	}
}
