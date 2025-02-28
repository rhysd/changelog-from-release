package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/v69/github"
)

type Config struct {
	Level        int
	Drafts       bool
	Prerelease   bool
	Contributors bool
	Ignore       *regexp.Regexp
	Extract      *regexp.Regexp
}

func (c *Config) filterReleases(rels []*github.RepositoryRelease) []*github.RepositoryRelease {
	i := 0
	for i < len(rels) {
		r := rels[i]
		t := r.GetTagName()
		if (c.Drafts || !r.GetDraft()) &&
			(c.Prerelease || !r.GetPrerelease()) &&
			(c.Ignore == nil || !c.Ignore.MatchString(t)) &&
			(c.Extract == nil || c.Extract.MatchString(t)) {
			i++
		} else {
			slog.Debug("Filtered release due to configuration", "release", r, "tag", t)
			rels = append(rels[:i], rels[i+1:]...)
		}
	}

	return rels
}

// GenerateChangeLog generates changelog text from given project data and configuration.
func GenerateChangeLog(c *Config, p *Project) ([]byte, error) {
	type ref struct {
		label string
		url   string
	}

	var out bytes.Buffer
	rels := c.filterReleases(p.Releases)
	heading := strings.Repeat("#", c.Level)
	url := p.RepoURL()

	slog.Debug("Start generating release notes", "url", url, "config", c)

	linker := NewReflinker(url)
	for _, l := range p.Autolinks {
		linker.AddExtRef(*l.KeyPrefix, *l.URLTemplate, *l.IsAlphanumeric)
	}

	knownUsers := make(map[string]bool)

	numRels := len(rels)
	refs := make([]ref, 0, numRels)
	for i, rel := range rels {
		prevTag := ""
		if i+1 < numRels {
			prevTag = rels[i+1].GetTagName()
		}

		title := strings.TrimSpace(rel.GetName())
		tag := rel.GetTagName()

		if tag == "" {
			return nil, fmt.Errorf(
				"release %q created at %s is not associated with any tag name. cannot determine a tag name for the release. did you forget setting tag name in the draft release?",
				strings.TrimSpace(rel.GetName()),
				rel.CreatedAt.Format(time.RFC3339),
			)
		}

		var created github.Timestamp
		if rel.GetDraft() {
			created = rel.GetCreatedAt()
		} else {
			created = rel.GetPublishedAt()
		}

		slog.Debug("Generating release", "name", title, "tag", tag, "created", created)

		var compareURL string
		if prevTag == "" {
			compareURL = fmt.Sprintf("%s/tree/%s", url, tag)
		} else {
			compareURL = fmt.Sprintf("%s/compare/%s...%s", url, prevTag, tag)
		}

		fmt.Fprintf(&out, "<a id=\"%s\"></a>\n", tag)

		if title == "" {
			title = tag
		} else if !strings.Contains(title, tag) {
			title = fmt.Sprintf("%s (%s)", title, tag)
		}

		pageURL := fmt.Sprintf("%s/releases/tag/%s", url, tag)
		date := created.Format(time.DateOnly)

		fmt.Fprintf(&out, "%s [%s](%s) - %s\n\n", heading, title, pageURL, date)
		fmt.Fprint(&out, linker.Link(strings.Replace(rel.GetBody(), "\r", "", -1)))

		generateContributorsSection(&out, linker, knownUsers, c)

		fmt.Fprintf(&out, "\n\n[Changes][%s]\n\n\n", tag)

		refs = append(refs, ref{tag, compareURL})

		slog.Debug("Generated release", "title", title, "page", pageURL, "date", date)
	}

	slog.Debug("Generate release links", "links", len(refs))
	for _, r := range refs {
		fmt.Fprintf(&out, "[%s]: %s\n", r.label, r.url)
	}

	fmt.Fprintf(&out, "\n<!-- Generated by https://github.com/rhysd/changelog-from-release %s -->\n", version)

	slog.Debug("Finish to generate release notes", "url", url)

	return out.Bytes(), nil
}

func generateContributorsSection(out *bytes.Buffer, linker *Reflinker, knownUsers map[string]bool, c *Config) {
	if !c.Contributors {
		return
	}

	var contributors []string
	home := linker.HomeURL()
	for _, n := range linker.Usernames() {
		if _, checked := knownUsers[n]; !checked {
			r, err := http.Head(fmt.Sprintf("%s/%s.png", home, n))
			// Verify user exists to avoid 404 on image load
			knownUsers[n] = err == nil && r.StatusCode == http.StatusOK
		}
		if knownUsers[n] {
			contributors = append(contributors, n)
		}
	}

	if len(contributors) == 0 {
		return
	}

	slog.Debug("Generating a contributors section", "contributors", contributors)
	fmt.Fprintf(out, "\n\n%s Contributors\n", strings.Repeat("#", c.Level+1))

	// Add profile image links
	for _, n := range contributors {
		u := url.QueryEscape(fmt.Sprintf("%s/%s.png", home, n))
		fmt.Fprintf(out, "\n<a href=\"%s/%s\"><img src=\"https://wsrv.nl/?url=%s&w=128&h=128&fit=cover&mask=circle\" width=\"64\" height=\"64\" alt=\"@%s\"></a>", home, n, u, n)
	}
}
