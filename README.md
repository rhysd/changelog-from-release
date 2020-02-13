Tiny ChangeLog Generator via GitHub Releases
============================================
[![macOS/Linux CI Status][travisci-badge]][travisci]
[![Windows CI Status][appveyor-badge]][appveyor]

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

## License

[the MIT License](LICENSE.txt)

[appveyor-badge]: https://ci.appveyor.com/api/projects/status/di0fr3r75afkrpkh?svg=true
[appveyor]: https://ci.appveyor.com/project/rhysd/changelog-from-release
[travisci-badge]: https://travis-ci.org/rhysd/changelog-from-release.svg?branch=master
[travisci]: https://travis-ci.org/rhysd/changelog-from-release
[releases]: https://github.com/rhysd/changelog-from-release/releases
