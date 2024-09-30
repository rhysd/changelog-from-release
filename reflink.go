package main

// Note: https://docs.github.com/en/get-started/writing-on-github/working-with-advanced-formatting/autolinked-references-and-urls

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
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

// Note: '_' is actually not boundary. But it's hard to check if the '_' is a part of italic/bold
// syntax.
// For example, _#123_ should be linked because '_'s are part of italic syntax. But _#123 and #123_
// should not be linked because '_'s are NOT part of italic syntax.
// Checking if the parent node is Italic/Bold or not does not help to solve this issue. For example,
// _foo_#1 should be linked. However #1 itself is not an italic text though the neighbor node is
// Italic.
// Fortunately this is very edge case. To keep our implementation simple, we compromise to treat '_'
// as a boundary. For example, _#1 and #1_ are linked incorrectly, but I believe they are OK for our
// use cases.
func isBoundary(b byte) bool {
	if '0' <= b && b <= '9' || 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' {
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

func (l *Reflinker) isBoundaryAt(idx int) bool {
	if idx < 0 || len(l.src) <= idx {
		return true
	}
	return isBoundary(l.src[idx])
}

func (l *Reflinker) lastIndexIssueRef(begin, end int) int {
	if !l.isBoundaryAt(begin - 1) {
		return -1 // Issue ref must follow a boundary (e.g. 'foo#bar')
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

	if !l.isBoundaryAt(end) {
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
	if !l.isBoundaryAt(begin - 1) {
		return -1 // e.g. foo@bar, _@foo (-@foo is ok)
	}

	// Note: Username may only contain alphanumeric characters or single hyphens, and cannot begin
	// or end with a hyphen: @foo-, @-foo
	// Note: '/' just after user name like @foo/ is not allowed

	if b := l.src[begin+1]; !isUserNameChar(b) || b == '-' {
		return -1
	}

	for i := 2; begin+i < end; i++ {
		b := l.src[begin+i]
		if isUserNameChar(b) {
			continue
		}
		if !isBoundary(b) || b == '/' || l.src[begin+i-1] == '-' {
			return -1
		}
		return begin + i
	}

	if l.src[end-1] == '-' {
		return -1
	}
	if end < len(l.src) {
		if b := l.src[end]; !isBoundary(b) || b == '/' {
			return -1
		}
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

const hashLen int = 40

func (l *Reflinker) linkCommitSHA(begin, end int) int {
	for i := 1; i < hashLen; i++ { // Since l.src[begin] was already checked, i starts from 1
		if begin+i >= end {
			return begin + i
		}
		b := l.src[begin+i]
		if '0' <= b && b <= '9' || 'a' <= b && b <= 'f' {
			continue
		}
		return begin + i
	}

	if l.isBoundaryAt(begin-1) && l.isBoundaryAt(begin+hashLen) {
		h := l.src[begin : begin+hashLen]
		l.links = append(l.links, refLink{
			start: begin,
			end:   begin + hashLen,
			text:  fmt.Sprintf("[`%s`](%s/commit/%s)", h[:10], l.repo, h),
		})
	}

	return begin + hashLen
}

// DetectLinksInText detects reference links in given markdown text and remembers them to replace all
// references later.
func (l *Reflinker) DetectLinksInText(t *ast.Text) {
	o := t.Segment.Start // start offset

	for o < t.Segment.Stop-1 { // `-1` means the last character is not checked
		s := l.src[o:t.Segment.Stop]
		i := bytes.IndexAny(s, "#@1234567890abcdef")
		if i < 0 || len(s)-1 <= i {
			return
		}
		switch s[i] {
		case '#':
			o = l.linkIssue(o+i, t.Segment.Stop)
		case '@':
			o = l.linkUser(o+i, t.Segment.Stop)
		default:
			// hex character [0-9a-f]
			o = l.linkCommitSHA(o+i, t.Segment.Stop)
		}
	}
}

var reGitHubCommitURL = regexp.MustCompile(`^https://github\.com/([^/]+/[^/]+)/commit/([[:xdigit:]]{7,})`)

func (l *Reflinker) detectGitHubURLInAutoLink(n *ast.AutoLink, src []byte) {
	start := 0
	p := n.PreviousSibling()
	if p != nil {
		t := p.(*ast.Text)
		if t == nil {
			return
		}
		start = t.Segment.Stop
	}

	label := n.Label(src)
	stop := start + len(label)
	if start >= len(src) || stop >= len(src) {
		return
	}
	if src[start] == '<' && stop+1 < len(src) && src[stop+1] == '>' {
		return
	}

	m := reGitHubCommitURL.FindSubmatch(label)
	if m == nil {
		return
	}

	slug := m[1]
	hash := m[2]
	if len(hash) > 10 {
		hash = hash[:10]
	}

	var replaced string
	if bytes.HasPrefix(label, []byte(l.repo)) {
		replaced = fmt.Sprintf("[`%s`](%s)", hash, label)
	} else {
		replaced = fmt.Sprintf("[`%s@%s`](%s)", slug, hash, label)
	}

	l.links = append(l.links, refLink{
		start: start,
		end:   stop,
		text:  replaced,
	})
}

// BuildLinkedText builds a markdown text where all references are replaced with links. The links were
// detected by DetectLinksInText() method calls.
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

// IsLinkDetected returns whether one or more links were detected.
func (l *Reflinker) IsLinkDetected() bool {
	return len(l.links) > 0
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

		switch n := n.(type) {
		case *ast.CodeSpan, *ast.Link:
			return ast.WalkSkipChildren, nil
		case *ast.AutoLink:
			l.detectGitHubURLInAutoLink(n, src)
			return ast.WalkSkipChildren, nil
		case *ast.Text:
			l.DetectLinksInText(n)
			return ast.WalkContinue, nil
		default:
			return ast.WalkContinue, nil
		}
	})

	if !l.IsLinkDetected() {
		return input
	}

	return l.BuildLinkedText()
}
