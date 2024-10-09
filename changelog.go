package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/v65/github"
)

type Config struct {
	Level   int
	Drafts  bool
	Ignore  *regexp.Regexp
	Extract *regexp.Regexp
}

func (c *Config) filterReleases(rels []*github.RepositoryRelease) []*github.RepositoryRelease {
	if c.Drafts && c.Ignore == nil && c.Extract == nil {
		return rels
	}

	i := 0
	for i < len(rels) {
		r := rels[i]
		t := r.GetTagName()
		if !c.Drafts && r.GetDraft() ||
			c.Ignore != nil && c.Ignore.MatchString(t) ||
			c.Extract != nil && !c.Extract.MatchString(t) {
			rels = append(rels[:i], rels[i+1:]...)
		} else {
			i++
		}
	}

	return rels
}

// ChangeLog is a struct to generate changelog output from given repository URL
type ChangeLog struct {
	out io.Writer
	cfg *Config
}

// Generate generates changelog text from given releases and outputs it to its writer
func (cl *ChangeLog) Generate(p *Project) error {
	type link struct {
		name string
		url  string
	}

	rels := cl.cfg.filterReleases(p.Releases)
	out := bufio.NewWriter(cl.out)
	heading := strings.Repeat("#", cl.cfg.Level)
	url := p.RepoURL()

	linker := NewReflinker(url)
	for _, l := range p.Autolinks {
		linker.AddExtRef(*l.KeyPrefix, *l.URLTemplate, *l.IsAlphanumeric)
	}

	numRels := len(rels)
	relLinks := make([]link, 0, numRels)
	for i, rel := range rels {
		prevTag := ""
		if i+1 < numRels {
			prevTag = rels[i+1].GetTagName()
		}

		title := strings.TrimSpace(rel.GetName())
		tag := rel.GetTagName()

		if tag == "" {
			return fmt.Errorf(
				"release %q created at %s is not associated with any tag name. cannot determine a tag name for the release. did you forget setting tag name in the draft release?",
				strings.TrimSpace(rel.GetName()),
				rel.CreatedAt.Format(time.RFC3339),
			)
		}

		var compareURL string
		if prevTag == "" {
			compareURL = fmt.Sprintf("%s/tree/%s", url, tag)
		} else {
			compareURL = fmt.Sprintf("%s/compare/%s...%s", url, prevTag, tag)
		}

		fmt.Fprintf(out, "<a name=\"%s\"></a>\n", tag)

		if title == "" {
			title = tag
		} else if !strings.Contains(title, tag) {
			title = fmt.Sprintf("%s (%s)", title, tag)
		}

		pageURL := fmt.Sprintf("%s/releases/tag/%s", url, tag)

		var created github.Timestamp
		if rel.GetDraft() {
			created = rel.GetCreatedAt()
		} else {
			created = rel.GetPublishedAt()
		}

		fmt.Fprintf(out, "%s [%s](%s) - %s\n\n", heading, title, pageURL, created.Format("02 Jan 2006"))
		fmt.Fprint(out, linker.Link(strings.Replace(rel.GetBody(), "\r", "", -1)))
		fmt.Fprintf(out, "\n\n[Changes][%s]\n\n\n", tag)

		relLinks = append(relLinks, link{tag, compareURL})
	}

	for _, l := range relLinks {
		fmt.Fprintf(out, "[%s]: %s\n", l.name, l.url)
	}

	fmt.Fprintf(out, "\n<!-- Generated by https://github.com/rhysd/changelog-from-release %s -->\n", version)

	return out.Flush()
}

// NewChangeLog creates a new ChangeLog instance
func NewChangeLog(w io.Writer, c *Config) *ChangeLog {
	return &ChangeLog{w, c}
}
