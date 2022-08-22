Generate Changelog from GitHub Releases
=======================================
[![CI][ci-badge]][ci]

`changelog-from-release` is a (too) small command line tool to generate changelog from
[GitHub Releases][gh-releases]. It fetches releases of the repository via GitHub API and generates
changelog in Markdown format.

For example, [CHANGELOG.md](./CHANGELOG.md) was generated from [the releases page][releases].

Other real-world examples:

- https://github.com/rhysd/actionlint/blob/main/CHANGELOG.md
- https://github.com/rhysd/hgrep/blob/main/CHANGELOG.md
- https://github.com/rhysd/git-brws/blob/master/CHANGELOG.md


## Installation

Download binary from [the releases page](https://github.com/rhysd/changelog-from-release/releases) or
build from sources with Go toolchain.

```
$ go install github.com/rhysd/changelog-from-release@latest
```


## Usage

Running `changelog-from-release` with no argument outputs a changelog text in Markdown format to
stdout. Please redirect the output to your changelog file.

```
$ cd /path/to/repo
$ changelog-from-release > CHANGELOG.md
$ cat CHANGELOG.md
```

Automation with [GitHub Actions][gh-actions] is also offered. Please read
[action's README](./action/README.md) for more details.

```yaml
- uses: rhysd/changelog-from-release/action@v2
  with:
    file: CHANGELOG.md
    github_token: ${{ secrets.GITHUB_TOKEN }}
```


## FAQ

### How to insert some templates at top/bottom of generated changelog?

Since `changelog-from-release` command just generates changelog history, you can insert your
favorite templates before/after redirecting the generated output to `CHANGELOG.md` file.

```sh
# Insert header
cat <<-EOS > CHANGELOG.md
Changelog
=========

This is a changelog for [my-project](https://github.com/owner/my-project

EOS

changelog-from-release >> CHANGELOG.md

# Insert footer
cat <<-EOS >> CHANGELOG.md

Releases on GitHub: https://github.com/owner/my-project/releases
EOS
```

If your shell supports `$()`, header and footer can be inserted once.

```sh
# Insert header
cat <<-EOS > CHANGELOG.md
Changelog
=========

This is a changelog for [my-project](https://github.com/owner/my-project

$(changelog-from-release)

Releases on GitHub: https://github.com/owner/my-project/releases
EOS
```

### How to update changelog before adding the release tag?

For example, how to include changes for v1.2.3 in `CHANGELOG.md` before creating a Git tag `v1.2.3`?

Please use [a release draft][gh-draft].

1. Create and save a new release note as draft
2. Run `changelog-from-release` with setting [a personal access token][pat] to `$GITHUB_TOKEN`
   environment variable
3. Commit the changelog
4. Create and push a new Git tag
5. Publish the release by clicking 'Publish release' button on GitHub

Setting a personal access token at 2. is mandatory since release drafts are private information.
API token associated with your account is necessary to fetch it.


## Reference auto linking

References in a release note are automatically converted to links by `changelog-from-release`.

- **Issue references** like `#123` are converted to links to the issue pages
- **User references** like `@rhysd` are converted to links to the user profile pages
- **Commit references** like `93e1af6ec49d23397baba466fba1e89cc8b6de39` are converted to linkes to the
  commit pages. To avoid false-positives, only full-length (40 characters) commit hashes are converted.

For example,

```markdown
Commit: 93e1af6ec49d23397baba466fba1e89cc8b6de39
Author: @rhysd
Issue:  #123
```

is converted to

```markdown
Commit: [`93e1af6ec4`](https://github.com/owner/repo/commit/93e1af6ec49d23397baba466fba1e89cc8b6de39)
Author: [@rhysd](https://github.com/rhysd)
Issue:  [#123](https://github.com/owner/repo/issues/123)
```


## Environment variables

### `GITHUB_API_BASE_URL`

For [GitHub Enterprise][ghe], please set `GITHUB_API_BASE_URL` environment variable to configure API
base URL.

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

[gh-releases]: https://docs.github.com/en/repositories/releasing-projects-on-github/about-releases
[releases]: https://github.com/rhysd/changelog-from-release/releases
[ci]: https://github.com/rhysd/changelog-from-release/actions?query=workflow%3ACI+branch%3Amaster
[ci-badge]: https://github.com/rhysd/changelog-from-release/workflows/CI/badge.svg?branch=master&event=push
[gh-actions]: https://github.com/features/actions
[ghe]: https://github.com/enterprise
[pat]: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token
[gh-draft]: https://docs.github.com/en/repositories/releasing-projects-on-github/managing-releases-in-a-repository
