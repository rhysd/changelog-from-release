package main

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-github/v45/github"
)

type link struct {
	name string
	url  string
}

// ChangeLog is a struct to generate changelog output from given repository URL
type ChangeLog struct {
	repoURL string
	out     io.Writer
	level   int
}

// Generate generates changelog text from given releases and outputs it to its writer
func (cl *ChangeLog) Generate(rels []*github.RepositoryRelease) error {
	out := bufio.NewWriter(cl.out)
	heading := strings.Repeat("#", cl.level)

	numRels := len(rels)
	relLinks := make([]link, 0, numRels)
	for i, rel := range rels {
		prevTag := ""
		if i+1 < numRels {
			prevTag = rels[i+1].GetTagName()
		}

		title := rel.GetName()
		tag := rel.GetTagName()

		if tag == "" {
			return fmt.Errorf(
				"release %q created at %s is not associated with any tag name. cannot determine a tag name for the release. did you forget setting tag name in the draft release?",
				rel.GetName(),
				rel.CreatedAt.Format(time.RFC3339),
			)
		}

		var compareURL string
		if prevTag == "" {
			compareURL = fmt.Sprintf("%s/tree/%s", cl.repoURL, tag)
		} else {
			compareURL = fmt.Sprintf("%s/compare/%s...%s", cl.repoURL, prevTag, tag)
		}

		fmt.Fprintf(out, "<a name=\"%s\"></a>\n", tag)

		if title == "" {
			title = tag
		} else if title != tag {
			title = fmt.Sprintf("%s (%s)", title, tag)
		}

		pageURL := fmt.Sprintf("%s/releases/tag/%s", cl.repoURL, tag)

		var created github.Timestamp
		if rel.GetDraft() {
			created = rel.GetCreatedAt()
		} else {
			created = rel.GetPublishedAt()
		}

		fmt.Fprintf(out, "%s [%s](%s) - %s\n\n", heading, title, pageURL, created.Format("02 Jan 2006"))
		fmt.Fprint(out, LinkRefs(strings.Replace(rel.GetBody(), "\r", "", -1), cl.repoURL))
		fmt.Fprintf(out, "\n\n[Changes][%s]\n\n\n", tag)

		relLinks = append(relLinks, link{tag, compareURL})
	}

	for _, link := range relLinks {
		fmt.Fprintf(out, "[%s]: %s\n", link.name, link.url)
	}

	fmt.Fprint(out, "\n <!-- Generated by https://github.com/rhysd/changelog-from-release -->\n")

	return out.Flush()
}

// NewChangeLog creates a new ChangeLog instance
func NewChangeLog(w io.Writer, u *url.URL, l int) *ChangeLog {
	// Strip credentials in the repository URL (#9)
	u.User = nil
	return &ChangeLog{strings.TrimSuffix(u.String(), ".git"), w, l}
}
