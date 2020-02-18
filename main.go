package main

import (
	"flag"
	"fmt"
	"os"
)

var version = "v1.1.3"

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
	ver := flag.Bool("version", false, "Output version to stdout")
	commit := flag.Bool("commit", false, "Create a new commit")
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

	if *commit {
		if err := git.CheckClean(); err != nil {
			fail(err)
		}
	}

	url, err := git.TrackingRemoteURL()
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
		fail(fmt.Errorf("No release was found at %s", url))
	}

	cl, err := NewChangeLog(git.root, url)
	if err != nil {
		fail(err)
	}

	if err := cl.Generate(rels); err != nil {
		fail(err)
	}

	if *commit {
		if err := git.CheckClean(); err == nil {
			fail(fmt.Errorf("Changelog file is up-to-date. Nothing to commit: %s", cl.file))
		}
		if err := git.Add(cl.file); err != nil {
			fail(err)
		}
		m := fmt.Sprintf("Update changelog for %s", rels[0].GetTagName())
		if err := git.Commit(m); err != nil {
			fail(err)
		}
	}
}
