package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v76/github"
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

func TestGenerateContributorsSection(t *testing.T) {
	l := NewReflinker("https://github.com/u/r")
	_ = l.Link("This is test for contributors section @rhysd @hello @world")
	c := &Config{Contributors: true, Level: 2}

	known := map[string]bool{"hello": true, "world": false}
	var b bytes.Buffer
	generateContributorsSection(&b, l, known, c)

	have := b.String()
	want := strings.Join([]string{
		"",
		"",
		"### Contributors",
		"",
		`<a href="https://github.com/hello"><img src="https://wsrv.nl/?url=https%3A%2F%2Fgithub.com%2Fhello.png&w=128&h=128&fit=cover&mask=circle" width="64" height="64" alt="@hello"></a>`,
		`<a href="https://github.com/rhysd"><img src="https://wsrv.nl/?url=https%3A%2F%2Fgithub.com%2Frhysd.png&w=128&h=128&fit=cover&mask=circle" width="64" height="64" alt="@rhysd"></a>`,
	}, "\n")
	if have != want {
		t.Fatal(cmp.Diff(have, want))
	}

	yes, ok := known["rhysd"]
	if !ok || !yes {
		t.Fatal("@rhysd is not registered as known user:", known)
	}
}

func TestDontGenerateContributorsSectionWithoutOption(t *testing.T) {
	l := NewReflinker("https://github.com/u/r")
	_ = l.Link("This is test for contributors section @rhysd @hello @world")
	c := &Config{}

	var b bytes.Buffer
	generateContributorsSection(&b, l, map[string]bool{"hello": true, "world": false}, c)
	out := b.String()
	if len(out) != 0 {
		t.Fatalf("contributors section should not be generated but got %q", out)
	}
}
