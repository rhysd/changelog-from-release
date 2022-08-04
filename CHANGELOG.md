<a name="v3.1.2"></a>
# [v3.1.2](https://github.com/rhysd/changelog-from-release/releases/tag/v3.1.2) - 04 Aug 2022

- Fixed `commit-summary-template` was not effective by renaming the input to `commit_summary_template`. Hyphens are not available in Docker actions. For example, when `v1.2.3` is released, the following step creates a commit with summary `chore(changelog): describe changes for "v1.2.3"`.
  ```yaml
  - uses: rhysd/changelog-from-release/action@v3
    with:
      file: CHANGELOG.md
      github_token: ${{ secrets.GITHUB_TOKEN }}
      commit_summary_template: 'chore(changelog): describe changes for %s'
  ```

[Changes][v3.1.2]


<a name="v3.1.1"></a>
# [v3.1.1](https://github.com/rhysd/changelog-from-release/releases/tag/v3.1.1) - 04 Aug 2022

- Use `git remote` instead of `git rev-parse` to retrieve a remote name of repository since `git rev-parse` sometimes returns an unexpected output for some reason. ([#6](https://github.com/rhysd/changelog-from-release/issues/6))
- Added `commit-summary-template` input to [the action](https://github.com/rhysd/changelog-from-release/tree/master/action) so that the commit message can be customized. The template is passed to the first argument of `printf` command. It must contain one `%s` placeholder which will be replaced with the tag name.
- Removed duplicate of command output from error messages on `git` command failure.
- Improved error messages when retrieving a URL of remote repository

[Changes][v3.1.1]


<a name="v3.1.0"></a>
# [v3.1.0](https://github.com/rhysd/changelog-from-release/releases/tag/v3.1.0) - 03 Aug 2022

- Link commit references in release note automatically. For example, `85a7d9028ed70bc81224cb126e29e070dcc0aa1c` is converted to ``[`85a7d9028e`](https://github.com/owner/repo/commit/85a7d9028ed70bc81224cb126e29e070dcc0aa1c)``. Note that only full-length (40 characters) commit hashes are linked to avoid false positives.
- Fix user references followed by `/` like `@foo/` are wrongly linked.
- Describe how reference auto linking works in [README.md](https://github.com/rhysd/changelog-from-release/blob/master/README.md).

[Changes][v3.1.0]


<a name="v3.0.0"></a>
# [v3.0.0](https://github.com/rhysd/changelog-from-release/releases/tag/v3.0.0) - 03 Aug 2022

- **BREAKING:** Labels at head of list items are no longer converted to bold text. For example, a list item starting with `- Fix:` was converted to `- **Fix:**`, but it is no longer converted from v3.0.0.
- Issue references like `#123` and user references like `@rhysd` are now automatically linked ([#3](https://github.com/rhysd/changelog-from-release/issues/3)). For example, when we have a release item as follows in release notes:
  ```markdown
  - Fixed something (thanks @rhysd, #1)
  ```
  `changelog-from-release` links the references as follows:
  ```markdown
  - Fixed something (thanks [@rhysd](https://github.com/rhysd), [#1](https://github.com/owner/repo/issues/1))
  ```
- Fixed `git@` and `ssh://` repository URLs were not converted to HTTPS URLs when the repository is hosted on GHE.
- Updated `google/go-github` dependency from v17 to v45.
- Removed `pkg/errors` dependency and used standard `fmt.Errorf` instead.

[Changes][v3.0.0]


<a name="v2.2.5"></a>
# [v2.2.5](https://github.com/rhysd/changelog-from-release/releases/tag/v2.2.5) - 02 Jun 2022

- Fix `changelog-from-release` command hangs when generating a changelog of repository which has more than 30 releases ([#8](https://github.com/rhysd/changelog-from-release/issues/8), [#10](https://github.com/rhysd/changelog-from-release/issues/10))

[Changes][v2.2.5]


<a name="v2.2.4"></a>
# [v2.2.4](https://github.com/rhysd/changelog-from-release/releases/tag/v2.2.4) - 12 May 2022

- Strip credentials in repository URLs ([#9](https://github.com/rhysd/changelog-from-release/issues/9)).
- Fix [the action](https://github.com/rhysd/changelog-from-release/tree/master/action) fails due to permission error on accessing a workspace directory.
- Update dependencies in `go.mod`.
- Use [GoReleaser](https://goreleaser.com/) to make release binaries.

[Changes][v2.2.4]


<a name="v2.2.3"></a>
# [v2.2.3](https://github.com/rhysd/changelog-from-release/releases/tag/v2.2.3) - 26 Sep 2021

- Improve: Introduce Go modules. Now this tool is installable via `go install`
- Improve: Better footer comment (thanks [@spl](https://github.com/spl), [#7](https://github.com/rhysd/changelog-from-release/issues/7))
- Improve: Build binaries with the latest Go toolchain v1.17
- Improve: Release `darwin/arm64` and `linux/arm64` binaries

[Changes][v2.2.3]


<a name="v2.2.2"></a>
# [v2.2.2](https://github.com/rhysd/changelog-from-release/releases/tag/v2.2.2) - 24 Feb 2021

- Fix: Rename `github-token` input to `github_token` since `-` is not available for input names ([#4](https://github.com/rhysd/changelog-from-release/issues/4))

[Changes][v2.2.2]


<a name="v2.2.0"></a>
# [v2.2.0](https://github.com/rhysd/changelog-from-release/releases/tag/v2.2.0) - 22 Feb 2020

- New: Support `$GITHUB_API_BASE_URL` environment variable to configure API endpoint for GitHub Enterprise

```sh
export GITHUB_API_BASE_URL=https://github.your-company.com/api/v3
GITHUB_TOKEN=abcabcabcabcabcabcabc changelog-from-release > CHANGELOG.md
```

[Changes][v2.2.0]


<a name="v2.1.0"></a>
# [v2.1.0](https://github.com/rhysd/changelog-from-release/releases/tag/v2.1.0) - 20 Feb 2020

- New: [Action](https://github.com/rhysd/changelog-from-release/tree/master/action) for [GitHub Actions]() was added. Updating your changelog file following the new release now can be automated easily.

Example workflow:

```yaml
name: Update changelog
on:
  release:
    types: [published]

jobs:
  changelog:
    name: Update changelog
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
        with:
          ref: master
      - uses: rhysd/changelog-from-release/action@v2
        with:
          file: CHANGELOG.md
          github-token: ${{ secrets.GITHUB_TOKEN }}
```


[Changes][v2.1.0]


<a name="v2.0.0"></a>
# [v2.0.0](https://github.com/rhysd/changelog-from-release/releases/tag/v2.0.0) - 19 Feb 2020

I decided to keep this tool as simple as possible. This release drops some features for simplification.

- Breaking: Instead of modifying `CHANGELOG.md`, this tool outputs a generated changelog to stdout. Please redirect the output to update your changelog file

```
$ changelog-from-release > CHANGELOG.md
```

- Breaking: Drop `-commit` flag. Please add and commit changes by yourself

```
$ changelog-from-release > CHANGELOG.md
$ git add CHANGELOG.md
$ git commit -m "Update changelog for $(git describe --tags)"
```

[Changes][v2.0.0]


<a name="v1.2.0"></a>
# [v1.2.0](https://github.com/rhysd/changelog-from-release/releases/tag/v1.2.0) - 18 Feb 2020

- New: `-commit` option was added to make a new commit for the changelog updates automatically

[Changes][v1.2.0]


<a name="v1.1.3"></a>
# [v1.1.3](https://github.com/rhysd/changelog-from-release/releases/tag/v1.1.3) - 13 Feb 2020

- Fix: Consider paging for getting releases from GitHub API
- Fix: Cause an error when no release found

[Changes][v1.1.3]


<a name="v1.1.2"></a>
# [v1.1.2](https://github.com/rhysd/changelog-from-release/releases/tag/v1.1.2) - 16 Nov 2019

- Fix: Codes in fences should not be modified

[Changes][v1.1.2]


<a name="v1.1.1"></a>
# [v1.1.1](https://github.com/rhysd/changelog-from-release/releases/tag/v1.1.1) - 14 Nov 2018

- Fix: Fix emphasizing item header with bold, not italic

[Changes][v1.1.1]


<a name="v1.1.0"></a>
# [v1.1.0](https://github.com/rhysd/changelog-from-release/releases/tag/v1.1.0) - 14 Nov 2018

- Improve: Emphasize list item headers like `- *Fix:* Fix something`

[Changes][v1.1.0]


<a name="v1.0.0"></a>
# [v1.0.0](https://github.com/rhysd/changelog-from-release/releases/tag/v1.0.0) - 10 Nov 2018

First release :tada:

[Changes][v1.0.0]


[v3.1.2]: https://github.com/rhysd/changelog-from-release/compare/v3.1.1...v3.1.2
[v3.1.1]: https://github.com/rhysd/changelog-from-release/compare/v3.1.0...v3.1.1
[v3.1.0]: https://github.com/rhysd/changelog-from-release/compare/v3.0.0...v3.1.0
[v3.0.0]: https://github.com/rhysd/changelog-from-release/compare/v2.2.5...v3.0.0
[v2.2.5]: https://github.com/rhysd/changelog-from-release/compare/v2.2.4...v2.2.5
[v2.2.4]: https://github.com/rhysd/changelog-from-release/compare/v2.2.3...v2.2.4
[v2.2.3]: https://github.com/rhysd/changelog-from-release/compare/v2.2.2...v2.2.3
[v2.2.2]: https://github.com/rhysd/changelog-from-release/compare/v2.2.0...v2.2.2
[v2.2.0]: https://github.com/rhysd/changelog-from-release/compare/v2.1.0...v2.2.0
[v2.1.0]: https://github.com/rhysd/changelog-from-release/compare/v2.0.0...v2.1.0
[v2.0.0]: https://github.com/rhysd/changelog-from-release/compare/v1.2.0...v2.0.0
[v1.2.0]: https://github.com/rhysd/changelog-from-release/compare/v1.1.3...v1.2.0
[v1.1.3]: https://github.com/rhysd/changelog-from-release/compare/v1.1.2...v1.1.3
[v1.1.2]: https://github.com/rhysd/changelog-from-release/compare/v1.1.1...v1.1.2
[v1.1.1]: https://github.com/rhysd/changelog-from-release/compare/v1.1.0...v1.1.1
[v1.1.0]: https://github.com/rhysd/changelog-from-release/compare/v1.0.0...v1.1.0
[v1.0.0]: https://github.com/rhysd/changelog-from-release/tree/v1.0.0

 <!-- Generated by https://github.com/rhysd/changelog-from-release -->
