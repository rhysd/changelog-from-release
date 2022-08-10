package main

import (
	"flag"
	"fmt"
	"os"
)

var version = "v3.1.4"

func usage() {
	fmt.Fprint(os.Stderr, "Usage: changelog-from-release [flags]\n\n")
	flag.PrintDefaults()
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	os.Exit(111)
}

func main() {
	flag.Usage = usage
	ver := flag.Bool("v", false, "Output version to stdout")
	flag.Parse()

	if *ver {
		fmt.Println(version)
		os.Exit(0)
	}

	if flag.NArg() != 0 {
		usage()
		os.Exit(111)
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

	cl := NewChangeLog(os.Stdout, url)

	if err := cl.Generate(rels); err != nil {
		fail(err)
	}
}
