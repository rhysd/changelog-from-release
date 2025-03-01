package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"regexp"
	"time"
)

const version = "v3.9.0"

func usage() {
	fmt.Fprint(os.Stderr, `Usage: changelog-from-release [FLAGS]

changelog-from-release is a (too) small command line tool to generate changelog
from GitHub Releases. This tool fetches releases of the repository via GitHub
API and generates the changelog to stdout in Markdown format. References like
#123 or @rhysd are automatically linked.

Please read the document for more details:

https://github.com/rhysd/changelog-from-release#readme

Flags:
`)
	flag.PrintDefaults()
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	os.Exit(111)
}

func regexFlag(v, f string) (*regexp.Regexp, error) {
	if v == "" {
		return nil, nil
	}
	r, err := regexp.Compile(v)
	if err != nil {
		return nil, fmt.Errorf("regular expression at %s is not valid: %w", f, err)
	}
	return r, nil
}

func remoteURL(config string) (*url.URL, error) {
	if config != "" {
		return ResolveRedirect(config)
	}

	git, err := NewGitForCwd()
	if err != nil {
		return nil, err
	}

	return git.FirstRemoteURL()
}

func fetchFromGitHub(u *url.URL) (*Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	gh, err := NewGitHub(u, ctx)
	if err != nil {
		return nil, err
	}

	return gh.Project()
}

func main() {
	flag.Usage = usage
	ver := flag.Bool("v", false, "Output version to stdout")
	heading := flag.Int("l", 1, "Heading level of each release section")
	drafts := flag.Bool("d", true, "Include draft releases")
	prerelease := flag.Bool("p", false, "Include pre-releases")
	ignore := flag.String("i", "", "Pattern to ignore release tags in regular expression")
	extract := flag.String("e", "", "Pattern to extract release tags in regular expression")
	remote := flag.String("r", "", "Remote repository URL to generate changelog")
	contrib := flag.Bool("c", false, "Generate contributors section")
	debug := flag.Bool("debug", false, "Enable debug log")
	flag.Parse()

	if *ver {
		fmt.Println(version)
		os.Exit(0)
	}
	if *debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	if *heading < 1 {
		fail(fmt.Errorf("heading level set by -l must be >=1 but %d is set", *heading))
	}

	reIgnore, err := regexFlag(*ignore, "-i")
	if err != nil {
		fail(err)
	}
	reExtract, err := regexFlag(*extract, "-e")
	if err != nil {
		fail(err)
	}
	cfg := &Config{
		Level:        *heading,
		Drafts:       *drafts,
		Prerelease:   *prerelease,
		Ignore:       reIgnore,
		Extract:      reExtract,
		Contributors: *contrib,
	}
	slog.Debug("Arguments parsed:", "config", cfg)

	if flag.NArg() != 0 {
		fail(fmt.Errorf("no argument is allowed but got %v", flag.Args()))
	}

	url, err := remoteURL(*remote)
	if err != nil {
		fail(err)
	}
	slog.Debug("Remote URL was resolved:", "config", *remote, "url", url)

	proj, err := fetchFromGitHub(url)
	if err != nil {
		fail(err)
	}
	slog.Debug("Fetched project via GitHub API:", "project", proj)

	gen, err := GenerateChangeLog(cfg, proj)
	if err != nil {
		fail(err)
	}

	slog.Debug("Write the generated output to stdout", "bytes", len(gen))
	if _, err := os.Stdout.Write(gen); err != nil {
		fail(fmt.Errorf("could not write the generated changelog to stdout: %w", err))
	}

	slog.Debug("Done")
}
