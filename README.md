Tiny ChangeLog Generator via GitHub Releases
============================================
[![CI][ci-badge]][ci]

`changelog-from-release` is a small command line tool to generate `CHANGELOG.md` at current directory.
It fetches releases of repository from GitHub API and generates `CHANGELOG.md`.

- From: [the releases page][releases]
- To: [CHANGELOG.md](./CHANGELOG.md)

Real-world examples:

- https://github.com/rhysd/notes-cli/blob/master/CHANGELOG.md
- https://github.com/rhysd/git-brws/blob/master/CHANGELOG.md

## Installation

Download binary from [release page](https://github.com/rhysd/changelog-from-release/releases) or
build from source with Go toolchain.

```
$ go get github.com/rhysd/changelog-from-release
```

## Usage

```
$ cd /path/to/repo
$ changelog-from-release
$ cat CHANGELOG.md
```

If you want to make a commit quickly for the changelog updates,

```
$ git commit -m "Update changelog for $(changelog-from-release -t)"
```

`-t` outputs the latest tag name to stdout.

Please see `changelog-from-release -h` for all options.

## License

[the MIT License](LICENSE.txt)

[releases]: https://github.com/rhysd/changelog-from-release/releases
[ci]: https://github.com/rhysd/changelog-from-release/actions?query=workflow%3ACI+branch%3Amaster
[ci-badge]: https://github.com/rhysd/changelog-from-release/workflows/CI/badge.svg?branch=master&event=push
