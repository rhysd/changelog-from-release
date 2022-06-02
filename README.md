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
$ go install github.com/rhysd/changelog-from-release
```

Note that `@latest` version specifier does not work at this moment.

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

## Environment variables

### `GITHUB_API_BASE_URL`

For GitHub Enterprise, please set `GITHUB_API_BASE_URL` environment variable to configure API base URL.

```sh
export GITHUB_API_BASE_URL=https://github.your-company.com/api/v3
```

### `GITHUB_TOKEN`

If `changelog-from-release` reported API rate limit exceeded or no permission to access the repository,
consider to specify [a personal access token][pat].

```sh
export GITHUB_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

## License

[the MIT License](LICENSE.txt)

[releases]: https://github.com/rhysd/changelog-from-release/releases
[ci]: https://github.com/rhysd/changelog-from-release/actions?query=workflow%3ACI+branch%3Amaster
[ci-badge]: https://github.com/rhysd/changelog-from-release/workflows/CI/badge.svg?branch=master&event=push
[gh-actions]: https://github.com/features/actions
[pat]: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token
