package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

var version = "v3.4.0"

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

func main() {
	flag.Usage = usage
	ver := flag.Bool("v", false, "Output version to stdout")
	heading := flag.Int("l", 1, "Heading level of each release section")
	ignore := flag.String("i", "", "Pattern to ignore release tags in regular expression")
	extract := flag.String("e", "", "Pattern to extract release tags in regular expression")
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

	git, err := NewGitForCwd()
	if err != nil {
		fail(err)
	}

	url, err := git.FirstRemoteURL()
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
	if len(rels) == 0 {
		fail(fmt.Errorf("no release was found at %s", url))
	}

	cl := NewChangeLog(os.Stdout, url, *heading, reIgnore, reExtract)

	if err := cl.Generate(rels); err != nil {
		fail(err)
	}
}
