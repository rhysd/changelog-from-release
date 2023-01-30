package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
)

const version = "v3.7.0"

func usage() {
	fmt.Fprint(os.Stderr, "Usage: changelog-from-release [flags]\n\n")
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

func main() {
	flag.Usage = usage
	ver := flag.Bool("v", false, "Output version to stdout")
	heading := flag.Int("l", 1, "Heading level of each release section")
	drafts := flag.Bool("d", true, "Include draft releases")
	ignore := flag.String("i", "", "Pattern to ignore release tags in regular expression")
	extract := flag.String("e", "", "Pattern to extract release tags in regular expression")
	remote := flag.String("r", "", "Remote repository URL to generate changelog")
	flag.Parse()

	if *ver {
		fmt.Println(version)
		os.Exit(0)
	}

	reIgnore, err := regexFlag(*ignore, "-i")
	if err != nil {
		fail(err)
	}
	reExtract, err := regexFlag(*extract, "-e")
	if err != nil {
		fail(err)
	}

	if flag.NArg() != 0 {
		usage()
		os.Exit(111)
	}

	if *heading < 1 {
		fail(fmt.Errorf("heading level set by -l must be >=1 but %d is set", *heading))
	}

	url, err := remoteURL(*remote)
	if err != nil {
		fail(err)
	}

	gh, err := GitHubFromURL(url)
	if err != nil {
		fail(err)
	}

	rels, err := gh.Releases()
	if err != nil {
		fail(err)
	}

	cl := NewChangeLog(os.Stdout, url, *heading, *drafts, reIgnore, reExtract)

	if err := cl.Generate(rels); err != nil {
		fail(err)
	}
}
