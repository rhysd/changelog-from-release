Tiny ChangeLog Generator via GitHub Releases
============================================
[![CI][ci-badge]][ci]

`changelog-from-release` is a (too) small command line tool to generate changelog.
It fetches releases of repository from GitHub API and generates changelog in Markdown format.

For example, [CHANGELOG.md](./CHANGELOG.md) is generated from [the releases page][releases].

Real-world examples:

- https://github.com/rhysd/notes-cli/blob/master/CHANGELOG.md
- https://github.com/rhysd/git-brws/blob/master/CHANGELOG.md

## Installation

Download binary from [the releases page](https://github.com/rhysd/changelog-from-release/releases) or
build from sources with Go toolchain.

```
$ go get github.com/rhysd/changelog-from-release
```

## Usage

Running `changelog-from-release` with no argument generates changelog in Markdown format and outputs
it to stdout. Please redirect the output to your changelog file for updating.

```
$ cd /path/to/repo
$ changelog-from-release > CHANGELOG.md
$ cat CHANGELOG.md
```

Automation with [GitHub Actions][gh-actions] is also offered.

```yaml
- uses: rhysd/changelog-from-release/action@v2
  with:
    file: CHANGELOG.md
    github_token: ${{ secrets.GITHUB_TOKEN }}
```

Please read [action's README](./action/README.md) for more details.

For GitHub Enterprise, please set `GITHUB_API_BASE_URL` environment variable.

```
export GITHUB_API_BASE_URL=https://github.your-company.com/api/v3
```

## License

[the MIT License](LICENSE.txt)

[releases]: https://github.com/rhysd/changelog-from-release/releases
[ci]: https://github.com/rhysd/changelog-from-release/actions?query=workflow%3ACI+branch%3Amaster
[ci-badge]: https://github.com/rhysd/changelog-from-release/workflows/CI/badge.svg?branch=master&event=push
[gh-actions]: https://github.com/features/actions
