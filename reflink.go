package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

type refLink struct {
	start int
	end   int
	text  string
}

func isBoundary(b byte) bool {
	if '0' <= b && b <= '9' || 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || b == '_' {
		return false
	}
	return true
}

func isUserNameChar(b byte) bool {
	return '0' <= b && b <= '9' || 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || b == '-'
}

// Reflinker detects all references in markdown text and replaces them with links.
type Reflinker struct {
	repo  string
	home  string
	src   []byte
	links []refLink
}

// NewReflinker creates Reflinker instance. repoURL is a repository URL of the service like
// https://github.com/user/repo.
func NewReflinker(repoURL string, src []byte) *Reflinker {
	u, err := url.Parse(repoURL)
	if err != nil {
		panic(err)
	}
	u.Path = ""

	return &Reflinker{
		repo:  repoURL,
		home:  u.String(),
		src:   src,
		links: nil,
	}
}

func (l *Reflinker) lastIndexIssueRef(begin, end int) int {
	if begin > 0 {
		// Ignore '_' to avoid hard edge case of GitHub's reference auto-linking behavior.
		//   _#1 → not linked (though we link it)
		//   _#1_ → linked
		//   _#1 foo_ → linked
		//   _foo_#1 → linked
		//   foo_#1 → not linked
		// We compromizes the first case so the behavior here is different from GitHub's behavior.
		// But it is very edge case and hard to handle correctly even if we use a Markdown parser.
		if b := l.src[begin-1]; b != '_' && !isBoundary(b) {
			return -1 // Issue ref must follow a boundary (e.g. 'foo#bar')
		}
	}

	for i := 1; begin+i < end; i++ {
		b := l.src[begin+i]
		if '0' <= b && b <= '9' {
			continue
		}
		if i == 1 || !isBoundary(b) {
			return -1
		}
		return begin + i
	}

	if end+1 < len(l.src) && !isBoundary(l.src[end+1]) {
		return -1
	}

	return end // The text ends with issue number
}

func (l *Reflinker) linkIssue(begin, end int) int {
	e := l.lastIndexIssueRef(begin, end)
	if e < 0 {
		return begin + 1
	}

	r := l.src[begin:e]
	l.links = append(l.links, refLink{
		start: begin,
		end:   e,
		text:  fmt.Sprintf("[%s](%s/issues/%s)", r, l.repo, r[1:]),
	})

	return e
}

func (l *Reflinker) lastIndexUserRef(begin, end int) int {
	if begin > 0 {
		// Ignore '_' to avoid hard edge case of GitHub's reference auto-linking behavior.
		//   _@x → not linked (though we link it)
		//   _@x_ → linked
		//   _@x foo_ → linked
		//   _foo_@x → linked
		//   foo_@x → not linked
		// We compromizes the first case so the behavior here is different from GitHub's behavior.
		// But it is very edge case and hard to handle correctly even if we use a Markdown parser.
		if b := l.src[begin-1]; b != '_' && !isBoundary(b) {
			return -1 // e.g. foo@bar, _@foo (-@foo is ok)
		}
	}

	// Username may only contain alphanumeric characters or single hyphens, and cannot begin or end
	// with a hyphen: @foo-, @-foo

	if b := l.src[begin+1]; !isUserNameChar(b) || b == '-' {
		return -1
	}

	for i := 2; begin+i < end; i++ {
		b := l.src[begin+i]
		if isUserNameChar(b) {
			continue
		}
		if !isBoundary(b) || l.src[begin+i-1] == '-' {
			return -1
		}
		return begin + i
	}

	if l.src[end-1] == '-' || end+1 < len(l.src) && !isBoundary(l.src[end+1]) {
		return -1
	}

	return end
}

func (l *Reflinker) linkUser(begin, end int) int {
	e := l.lastIndexUserRef(begin, end)
	if e < 0 {
		return begin + 1
	}

	u := l.src[begin:e]
	l.links = append(l.links, refLink{
		start: begin,
		end:   e,
		text:  fmt.Sprintf("[%s](%s/%s)", u, l.home, u[1:]),
	})

	return e
}

// DetectLinks detects reference links in given markdown text and remembers them to replace all
// references later.
func (l *Reflinker) DetectLinks(t *ast.Text) {
	o := t.Segment.Start // start offset

	for o < t.Segment.Stop-1 { // `-1` means the last character is not checked
		s := l.src[o:t.Segment.Stop]
		i := bytes.IndexAny(s, "#@")
		if i < 0 || len(s)-1 <= i {
			return
		}
		switch s[i] {
		case '#':
			o = l.linkIssue(o+i, t.Segment.Stop)
		case '@':
			o = l.linkUser(o+i, t.Segment.Stop)
		}
	}
}

// BuildLinkedText builds a markdown text where all references are replaced with links. The links were
// detected by DetectLinks() method calls.
func (l *Reflinker) BuildLinkedText() string {
	if len(l.links) == 0 {
		return string(l.src)
	}

	var b strings.Builder
	i := 0
	for _, r := range l.links {
		b.Write(l.src[i:r.start])
		b.WriteString(r.text)
		i = r.end
	}
	b.Write(l.src[i:])
	return b.String()
}

// LinkRefs replaces all references in the given markdown text with actual links.
func LinkRefs(input string, repoURL string) string {
	src := []byte(input)
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	t := md.Parser().Parse(text.NewReader(src))
	l := NewReflinker(repoURL, src)

	ast.Walk(t, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if n, ok := n.(*ast.Text); ok {
			if _, ok := n.Parent().(*ast.CodeSpan); !ok {
				l.DetectLinks(n)
			}
		}

		return ast.WalkContinue, nil
	})

	return l.BuildLinkedText()
}
