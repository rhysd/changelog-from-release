ChangeLog Generator via GitHub Releases
=======================================

This is a small command line tool to generate `CHANGELOG.md` at current directory.
It fetches releases of repoisitory of current directory and generates `CHANGELOG.md` with them.

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

[MIT License](LICENSE.txt)
