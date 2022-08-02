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

func lastIndexIssueRef(src []byte, begin int) int {
	if begin > 0 && !isBoundary(src[begin-1]) {
		return -1 // Issue ref must follow a boundary (e.g. 'foo#bar')
	}

	for i := 1; begin+i < len(src); i++ {
		b := src[begin+i]
		if '0' <= b && b <= '9' {
			continue
		}
		if i == 1 || !isBoundary(b) {
			return -1
		}
		return begin + i
	}

	return len(src) // The text ends with issue number
}

func isUserNameChar(b byte) bool {
	return '0' <= b && b <= '9' || 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || b == '-'
}

func lastIndexUserRef(src []byte, begin int) int {
	if begin > 0 && !isBoundary(src[begin-1]) {
		return -1 // e.g. foo@bar, _@foo (-@foo is ok)
	}

	// Username may only contain alphanumeric characters or single hyphens, and cannot begin or end
	// with a hyphen: @foo-, @-foo

	if b := src[begin+1]; !isUserNameChar(b) || b == '-' {
		return -1
	}

	for i := 2; begin+i < len(src); i++ {
		b := src[begin+i]
		if isUserNameChar(b) {
			continue
		}
		if !isBoundary(b) || src[begin+i-1] == '-' {
			return -1
		}
		return begin + i
	}

	if src[len(src)-1] == '-' {
		return -1
	}

	return len(src)
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

func (l *Reflinker) linkIssue(src []byte, begin, offset int) int {
	e := lastIndexIssueRef(src, begin)
	if e < 0 {
		return begin + 1
	}

	r := src[begin:e]
	l.links = append(l.links, refLink{
		start: offset + begin,
		end:   offset + e,
		text:  fmt.Sprintf("[%s](%s/issues/%s)", r, l.repo, r[1:]),
	})

	return e
}

func (l *Reflinker) linkUser(src []byte, begin, offset int) int {
	e := lastIndexUserRef(src, begin)
	if e < 0 {
		return begin + 1
	}

	u := src[begin:e]
	l.links = append(l.links, refLink{
		start: offset + begin,
		end:   offset + e,
		text:  fmt.Sprintf("[%s](%s/%s)", u, l.home, u[1:]),
	})

	return e
}

// Link detects reference links in given markdown text and remembers them to replace all references
// later.
func (l *Reflinker) Link(t *ast.Text) {
	s := l.src[t.Segment.Start:t.Segment.Stop]
	o := 0
	for len(s) > 1 {
		b := bytes.IndexAny(s, "#@")
		if b < 0 || b == len(s)-1 {
			return
		}
		switch s[b] {
		case '#':
			i := l.linkIssue(s, b, t.Segment.Start+o)
			s = s[i:]
			o += i
		case '@':
			i := l.linkUser(s, b, t.Segment.Start+o)
			s = s[i:]
			o += i
		}
	}
}

// Build builds a markdown text where all references are replaced with links. The links were
// detected by Link() method calls.
func (l *Reflinker) Build() string {
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
				l.Link(n)
			}
		}

		return ast.WalkContinue, nil
	})

	return l.Build()
}
