<a id="v3.9.1"></a>
# [v3.9.1](https://github.com/rhysd/changelog-from-release/releases/tag/v3.9.1) - 2025-10-19

- Fix extracting the latest version from `git tag` output in the action. ([#39](https://github.com/rhysd/changelog-from-release/issues/39), thanks [@wu-clan](https://github.com/wu-clan))
- Update Go dependencies dropping support for Go 1.23 and earlier.

[Changes][v3.9.1]


<a id="v3.9.0"></a>
# [v3.9.0](https://github.com/rhysd/changelog-from-release/releases/tag/v3.9.0) - 2025-02-28

- Add `-c` option to generate "Contributors" section like below. Here is the [example of actionlint repository](https://gist.github.com/rhysd/fabe27f6c1610cd4fd84e85c134102d8) enabling the option. There are some notes on this option. Please read the [document](https://github.com/rhysd/changelog-from-release?tab=readme-ov-file#can-i-generate-contributors-section) for more details. ([#33](https://github.com/rhysd/changelog-from-release/issues/33), thanks [@yottahmd](https://github.com/yottahmd))
  ```markdown
  ## Contributors

  <a href="https://github.com/rhysd"><img src="https://wsrv.nl/?url=https%3A%2F%2Fgithub.com%2Frhysd.png&w=128&h=128&fit=cover&mask=circle" width="64" height="64" alt="@rhysd"></a>
  <a href="https://github.com/yottahmd"><img src="https://wsrv.nl/?url=https%3A%2F%2Fgithub.com%2Fyottahmd.png&w=128&h=128&fit=cover&mask=circle" width="64" height="64" alt="@yottahmd"></a>
  ```
- Fix crash when some external reference prefix contains a regex meta character.

[Changes][v3.9.0]


<a id="v3.8.1"></a>
# [v3.8.1](https://github.com/rhysd/changelog-from-release/releases/tag/v3.8.1) - 2024-12-02

- Fix resolving redirect of repository URL does not work when the repository is hosted on a GitHub Enterprise server. changelog-from-release tries to resolve a repository URL because the repository may have been renamed. However GitHub Enterprise server redirects a repository URL to its login page when the request is not authenticated. Now changelog-from-release detects a login page URL and stops resolving the redirect.

[Changes][v3.8.1]


<a id="v3.8.0"></a>
# [v3.8.0](https://github.com/rhysd/changelog-from-release/releases/tag/v3.8.0) - 2024-10-15

- Support more auto link references following the [GitHub's official spec](https://docs.github.com/en/get-started/writing-on-github/working-with-advanced-formatting/autolinked-references-and-urls). Please read the ['Reference auto linking' section](https://github.com/rhysd/changelog-from-release?tab=readme-ov-file#reference-auto-linking) for more details. Some examples:
  - Commit URL: [`50b11ed2bd`](https://github.com/rhysd/changelog-from-release/commit/50b11ed2bd8ce4efbc62770c19cc4e6eab74c7f6), [rhysd/actionlint@`1ba25a77e1`](https://github.com/rhysd/actionlint/commit/1ba25a77e1bece39b3b722a235d1e31ef84cd329)
  - Issue/PR URL: [#1](https://github.com/rhysd/changelog-from-release/issues/1), [rhysd/actionlint#453](https://github.com/rhysd/actionlint/pull/453)
  - `GH-` reference link: [GH-1](https://github.com/rhysd/changelog-from-release/issues/1)
- Support [the custom autolinks](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/managing-repository-settings/configuring-autolinks-to-reference-external-resources). Please read the ['Custom autolink' section](https://github.com/rhysd/changelog-from-release?tab=readme-ov-file#custom-autolink) for more details.
  - GitHub API requires "Administration" repository permissions (read). You need to set your personal access token to the `GITHUB_TOKEN` environment variable with proper permission.
- Fix pre-releases are included in a generated changelog. ([#31](https://github.com/rhysd/changelog-from-release/issues/31))
- Add `-p` flag to include pre-releases in a generated changelog.
- Prefer the ISO-standard date format `YYYY-MM-DD`. ([#27](https://github.com/rhysd/changelog-from-release/issues/27))
  - This follows the recommendation by [keep a changelog](https://keepachangelog.com/en/1.1.0/).
  - Previous format (e.g. `02 Jan 2006`) was not good because not all release notes are written in English.
- Use `id` attribute instead of `name` attribute for `<a>` elements in a generated changelog because the `name` attribute is deprecated.
- Fix releases are not correctly filtered when both `-d` and `-i`/`-e` are specified.
- Set 120 seconds timeout to GitHub API requests.
- Send requests to GitHub in parallel to fetch repository data faster.
- Add `-debug` flag and debug log. When something went wrong, enabling debug log helps to analyze what happened.
- Include checksums for each released archives in the release assets as `changelog-from-release_{version}_checksums.txt`.
- Require Go 1.22 or later for build.
- Update `go-github` dependency from v58 to v66.

[Changes][v3.8.0]


<a id="v3.7.2"></a>
# [v3.7.2](https://github.com/rhysd/changelog-from-release/releases/tag/v3.7.2) - 2024-01-26

- Fix getting a tag from `$GITHUB_EVENT_PATH` environment variable in action ([#23](https://github.com/rhysd/changelog-from-release/issues/23), thanks [@linde12](https://github.com/linde12))
- Fix only `github.com` is allowed as host name for GHE environment
- Do not download files on Git LFS when updating changelog file in action ([#24](https://github.com/rhysd/changelog-from-release/issues/24))
- Update Go module dependencies including `go-github` v58

[Changes][v3.7.2]


<a id="v3.7.1"></a>
# [v3.7.1](https://github.com/rhysd/changelog-from-release/releases/tag/v3.7.1) - 2023-04-12

- Ensure a trailing `/` in the API base URL set in `GITHUB_API_BASE_URL` environment variable
- Show a diff of updated changelog in action output instead of printing an entire changelog
- Update `google/go-github` from v40 to v45
- Improve help description of `-r` option

[Changes][v3.7.1]


<a id="v3.7.0"></a>
# [v3.7.0](https://github.com/rhysd/changelog-from-release/releases/tag/v3.7.0) - 2023-01-29

- Add `-r` option to specify a remote URL of repository.
  ```sh
  # Generate changelog for rhysd/changelog-from-release
  changelog-from-release -r 'https://github.com/rhysd/changelog-from-release'
  ```
- Fix repeating Git tag name in a release heading when a release title already includes it ([#20](https://github.com/rhysd/changelog-from-release/issues/20)). For example, when a release title is `v1.2.3 with some features` and its Git tag is `v1.2.3`, the generated heading is:
  - until v3.7.0: `v1.2.3 with some features (v1.2.3)`
  - from v3.7.0: `v1.2.3 with some features`
- Ensure spaces are trimmed from release title and release name.

[Changes][v3.7.0]


<a id="v3.6.1"></a>
# [v3.6.1](https://github.com/rhysd/changelog-from-release/releases/tag/v3.6.1) - 2023-01-16

- Fix 404 response is not handled when trying to resolve private renamed repositories. ([#19](https://github.com/rhysd/changelog-from-release/issues/19))
  - For private repositories, repository rename is not resolved because GitHub always returns 404 even if an authentication token is set. Please ensure the Git remote URL in your local repository is up-to-date when running `changelog-from-release` command in this case.

[Changes][v3.6.1]


<a id="v3.6.0"></a>
# [v3.6.0](https://github.com/rhysd/changelog-from-release/releases/tag/v3.6.0) - 2023-01-13

- If you enable [protected-branch](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/about-protected-branches), `rhysd/changelog-from-release/action` action cannot push a commit directly to the branch. Instead, use `pull_request` input to create a pull request to update the changelog. ([#17](https://github.com/rhysd/changelog-from-release/issues/17))
  ```yaml
  - uses: rhysd/changelog-from-release/action@v3
    with:
      file: CHANGELOG.md
      github_token: ${{ secrets.GITHUB_TOKEN }}
      pull_request: true
  ```

[Changes][v3.6.0]


<a id="v3.5.2"></a>
# [v3.5.2](https://github.com/rhysd/changelog-from-release/releases/tag/v3.5.2) - 2023-01-05

- Check and follow redirects for Git remote URLs. This check is necessary to resolve renamed old repositories correctly. ([#16](https://github.com/rhysd/changelog-from-release/issues/16))
- Avoid unnecessary memory allocation when no reference link is included in changelog.
- Update dependencies to include `golang.org/x/*` packages which were newly managed as Go modules.

[Changes][v3.5.2]


<a id="v3.5.1"></a>
# [v3.5.1](https://github.com/rhysd/changelog-from-release/releases/tag/v3.5.1) - 2022-12-12

- Add `-d` option to include/exclude drafts in generated changelog. (thanks [@paescuj](https://github.com/paescuj), [#15](https://github.com/rhysd/changelog-from-release/issues/15))
  ```sh
  # Exclude drafts from the output
  changelog-from-release -d=false > CHANGELOG.md
  ```
- Include version of `changelog-from-release` in the footer of generated output
- Remove a single space which were prepended to a footer line
- Do not fail when no release is found since an empty changelog is a good start point of development

[Changes][v3.5.1]


<a id="v3.4.0"></a>
# [v3.4.0](https://github.com/rhysd/changelog-from-release/releases/tag/v3.4.0) - 2022-08-27

- Add `-i` option to ignore release tags by regular expression pattern. For example, if your project has `nightly` tag release for nightly builds, it can be excluded as follows:
  ```sh
  changelog-from-release -i '^nightly$' > CHANGELOG.md
  ```
- Add `-e` option to extract release tags by regular expression pattern. For example, if your project uses `v{major}.{minor}.{patch}` format for release tags, they can be extracted as follows:
  ```sh
  changelog-from-release -e '^v\d+\.\d+\.\d+$' > CHANGELOG.md
  ```
- Allow multiple drafts in releases. For including draft releases, see [the FAQ](https://github.com/rhysd/changelog-from-release#how-to-update-changelog-before-adding-the-release-tag) for more details.

[Changes][v3.4.0]


<a id="v3.3.0"></a>
# [v3.3.0](https://github.com/rhysd/changelog-from-release/releases/tag/v3.3.0) - 2022-08-23

- Add `args` input to [the action](https://github.com/rhysd/changelog-from-release/tree/master/action) to define command line arguments passed to `changelog-from-release` command.
  ```yaml
  - uses: rhysd/changelog-from-release/action@v3
    with:
      file: CHANGELOG.md
      github_token: ${{ secrets.GITHUB_TOKEN }}
      # Pass `-l 2` to use `##` instead of `#` for each release section
      args: -l 2
  ```
- Add `header` and `footer` inputs to [the action](https://github.com/rhysd/changelog-from-release/tree/master/action) to insert templates before/after the generated changelog. The following step inserts the header and the footer.
  ```yaml
  - uses: rhysd/changelog-from-release/action@v3
    with:
      file: CHANGELOG.md
      github_token: ${{ secrets.GITHUB_TOKEN }}
      args: -l 2
      header: |
        Changelog
        =========

        This is header.
      footer: |-

        This is footer.
  ```
- Report an error when the release is not associated with any Git tags. This can happen when the release is a draft.
- Fix release date is broken when the release is a draft. Instead of published date, created date is used in the case.
- Add [FAQ section](https://github.com/rhysd/changelog-from-release#faq) to readme document. Currently two topics are described.
  - [How to update changelog before adding the release tag?](https://github.com/rhysd/changelog-from-release#how-to-update-changelog-before-adding-the-release-tag)
  - [How to insert some templates at top/bottom of generated changelog?](https://github.com/rhysd/changelog-from-release#how-to-insert-some-templates-at-topbottom-of-generated-changelog)

[Changes][v3.3.0]


<a id="v3.2.0"></a>
# [v3.2.0](https://github.com/rhysd/changelog-from-release/releases/tag/v3.2.0) - 2022-08-22

- Add `-l` option to set heading level of each release sections. For example, `-l 2` uses `##` instead of `#` for each release sections.
- Use Go 1.19 to build release binaries.
- Go module path was changed from `github.com/rhysd/changelog-from-release` to `github.com/rhysd/changelog-from-release/v3` since `go install` without version specifier does not work in recent Go toolchain. ([#14](https://github.com/rhysd/changelog-from-release/issues/14))

[Changes][v3.2.0]


<a id="v3.1.4"></a>
# [v3.1.4](https://github.com/rhysd/changelog-from-release/releases/tag/v3.1.4) - 2022-08-10

- In previous release, references in link texts were fixed. But the fix was not perfect. Nested text node in a link node was still linked incorrectly and this release fixed the bug. For instance, `@foo` in `[_@foo_](...)` should not be linked where the text is in italic node in link node.

[Changes][v3.1.4]


<a id="v3.1.3"></a>
# [v3.1.3](https://github.com/rhysd/changelog-from-release/releases/tag/v3.1.3) - 2022-08-08

- Do not references in link texts. Previously references in link texts like `[#1](...)` were linked to `[[#1](https://github.com/owner/repo/issues/1)](...)`. Now they are not linked and left as-is ([#12](https://github.com/rhysd/changelog-from-release/issues/12)).


[Changes][v3.1.3]


<a id="v3.1.2"></a>
# [v3.1.2](https://github.com/rhysd/changelog-from-release/releases/tag/v3.1.2) - 2022-08-04

- Fixed `commit-summary-template` was not effective by renaming the input to `commit_summary_template`. Hyphens are not available in Docker actions. For example, when `v1.2.3` is released, the following step creates a commit with summary `chore(changelog): describe changes for "v1.2.3"`.
  ```yaml
  - uses: rhysd/changelog-from-release/action@v3
    with:
      file: CHANGELOG.md
      github_token: ${{ secrets.GITHUB_TOKEN }}
      commit_summary_template: 'chore(changelog): describe changes for %s'
  ```

[Changes][v3.1.2]


<a id="v3.1.1"></a>
# [v3.1.1](https://github.com/rhysd/changelog-from-release/releases/tag/v3.1.1) - 2022-08-04

- Use `git remote` instead of `git rev-parse` to retrieve a remote name of repository since `git rev-parse` sometimes returns an unexpected output for some reason. ([#6](https://github.com/rhysd/changelog-from-release/issues/6))
- Added `commit-summary-template` input to [the action](https://github.com/rhysd/changelog-from-release/tree/master/action) so that the commit message can be customized. The template is passed to the first argument of `printf` command. It must contain one `%s` placeholder which will be replaced with the tag name.
- Removed duplicate of command output from error messages on `git` command failure.
- Improved error messages when retrieving a URL of remote repository

[Changes][v3.1.1]


<a id="v3.1.0"></a>
# [v3.1.0](https://github.com/rhysd/changelog-from-release/releases/tag/v3.1.0) - 2022-08-03

- Link commit references in release note automatically. For example, `85a7d9028ed70bc81224cb126e29e070dcc0aa1c` is converted to ``[`85a7d9028e`](https://github.com/owner/repo/commit/85a7d9028ed70bc81224cb126e29e070dcc0aa1c)``. Note that only full-length (40 characters) commit hashes are linked to avoid false positives.
- Fix user references followed by `/` like `@foo/` are wrongly linked.
- Describe how reference auto linking works in [README.md](https://github.com/rhysd/changelog-from-release/blob/master/README.md).

[Changes][v3.1.0]


<a id="v3.0.0"></a>
# [v3.0.0](https://github.com/rhysd/changelog-from-release/releases/tag/v3.0.0) - 2022-08-03

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


<a id="v2.2.5"></a>
# [v2.2.5](https://github.com/rhysd/changelog-from-release/releases/tag/v2.2.5) - 2022-06-02

- Fix `changelog-from-release` command hangs when generating a changelog of repository which has more than 30 releases ([#8](https://github.com/rhysd/changelog-from-release/issues/8), [#10](https://github.com/rhysd/changelog-from-release/issues/10))

[Changes][v2.2.5]


<a id="v2.2.4"></a>
# [v2.2.4](https://github.com/rhysd/changelog-from-release/releases/tag/v2.2.4) - 2022-05-12

- Strip credentials in repository URLs ([#9](https://github.com/rhysd/changelog-from-release/issues/9)).
- Fix [the action](https://github.com/rhysd/changelog-from-release/tree/master/action) fails due to permission error on accessing a workspace directory.
- Update dependencies in `go.mod`.
- Use [GoReleaser](https://goreleaser.com/) to make release binaries.

[Changes][v2.2.4]


<a id="v2.2.3"></a>
# [v2.2.3](https://github.com/rhysd/changelog-from-release/releases/tag/v2.2.3) - 2021-09-26

- Improve: Introduce Go modules. Now this tool is installable via `go install`
- Improve: Better footer comment (thanks [@spl](https://github.com/spl), [#7](https://github.com/rhysd/changelog-from-release/issues/7))
- Improve: Build binaries with the latest Go toolchain v1.17
- Improve: Release `darwin/arm64` and `linux/arm64` binaries

[Changes][v2.2.3]


<a id="v2.2.2"></a>
# [v2.2.2](https://github.com/rhysd/changelog-from-release/releases/tag/v2.2.2) - 2021-02-24

- Fix: Rename `github-token` input to `github_token` since `-` is not available for input names ([#4](https://github.com/rhysd/changelog-from-release/issues/4))

[Changes][v2.2.2]


<a id="v2.2.0"></a>
# [v2.2.0](https://github.com/rhysd/changelog-from-release/releases/tag/v2.2.0) - 2020-02-22

- New: Support `$GITHUB_API_BASE_URL` environment variable to configure API endpoint for GitHub Enterprise

```sh
export GITHUB_API_BASE_URL=https://github.your-company.com/api/v3
GITHUB_TOKEN=abcabcabcabcabcabcabc changelog-from-release > CHANGELOG.md
```

[Changes][v2.2.0]


<a id="v2.1.0"></a>
# [v2.1.0](https://github.com/rhysd/changelog-from-release/releases/tag/v2.1.0) - 2020-02-20

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


<a id="v2.0.0"></a>
# [v2.0.0](https://github.com/rhysd/changelog-from-release/releases/tag/v2.0.0) - 2020-02-19

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


<a id="v1.2.0"></a>
# [v1.2.0](https://github.com/rhysd/changelog-from-release/releases/tag/v1.2.0) - 2020-02-18

- New: `-commit` option was added to make a new commit for the changelog updates automatically

[Changes][v1.2.0]


<a id="v1.1.3"></a>
# [v1.1.3](https://github.com/rhysd/changelog-from-release/releases/tag/v1.1.3) - 2020-02-13

- Fix: Consider paging for getting releases from GitHub API
- Fix: Cause an error when no release found

[Changes][v1.1.3]


<a id="v1.1.2"></a>
# [v1.1.2](https://github.com/rhysd/changelog-from-release/releases/tag/v1.1.2) - 2019-11-16

- Fix: Codes in fences should not be modified

[Changes][v1.1.2]


<a id="v1.1.1"></a>
# [v1.1.1](https://github.com/rhysd/changelog-from-release/releases/tag/v1.1.1) - 2018-11-14

- Fix: Fix emphasizing item header with bold, not italic

[Changes][v1.1.1]


<a id="v1.1.0"></a>
# [v1.1.0](https://github.com/rhysd/changelog-from-release/releases/tag/v1.1.0) - 2018-11-14

- Improve: Emphasize list item headers like `- *Fix:* Fix something`

[Changes][v1.1.0]


<a id="v1.0.0"></a>
# [v1.0.0](https://github.com/rhysd/changelog-from-release/releases/tag/v1.0.0) - 2018-11-10

First release :tada:

[Changes][v1.0.0]


[v3.9.1]: https://github.com/rhysd/changelog-from-release/compare/v3.9.0...v3.9.1
[v3.9.0]: https://github.com/rhysd/changelog-from-release/compare/v3.8.1...v3.9.0
[v3.8.1]: https://github.com/rhysd/changelog-from-release/compare/v3.8.0...v3.8.1
[v3.8.0]: https://github.com/rhysd/changelog-from-release/compare/v3.7.2...v3.8.0
[v3.7.2]: https://github.com/rhysd/changelog-from-release/compare/v3.7.1...v3.7.2
[v3.7.1]: https://github.com/rhysd/changelog-from-release/compare/v3.7.0...v3.7.1
[v3.7.0]: https://github.com/rhysd/changelog-from-release/compare/v3.6.1...v3.7.0
[v3.6.1]: https://github.com/rhysd/changelog-from-release/compare/v3.6.0...v3.6.1
[v3.6.0]: https://github.com/rhysd/changelog-from-release/compare/v3.5.2...v3.6.0
[v3.5.2]: https://github.com/rhysd/changelog-from-release/compare/v3.5.1...v3.5.2
[v3.5.1]: https://github.com/rhysd/changelog-from-release/compare/v3.4.0...v3.5.1
[v3.4.0]: https://github.com/rhysd/changelog-from-release/compare/v3.3.0...v3.4.0
[v3.3.0]: https://github.com/rhysd/changelog-from-release/compare/v3.2.0...v3.3.0
[v3.2.0]: https://github.com/rhysd/changelog-from-release/compare/v3.1.4...v3.2.0
[v3.1.4]: https://github.com/rhysd/changelog-from-release/compare/v3.1.3...v3.1.4
[v3.1.3]: https://github.com/rhysd/changelog-from-release/compare/v3.1.2...v3.1.3
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

<!-- Generated by https://github.com/rhysd/changelog-from-release v3.9.1 -->
