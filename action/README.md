Automate changelog updates with [GitHub Actions][gh-actions]
============================================================
![Release version][release-badge]

This directory is an action for [GitHub Actions][gh-actions] to automate changelog updates.
It updates the changelog file, makes a commit and pushes it to the remote.

```yaml
name: On release published
on:
  release:
    types: [published]

jobs:
  changelog:
    name: Update changelog
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          ref: main
      - uses: rhysd/changelog-from-release/action@v3
        with:
          file: CHANGELOG.md
          github_token: ${{ secrets.GITHUB_TOKEN }}
```

`file` is a file path for changelog and `github_token` is a GitHub API token to retrieve releases.
Some other inputs are offered for customizing the behavior. Please read [action.yml](./action.yml)
for more details.

Note that `actions/checkout@v3` does not fetch branches by default. In above example, `ref: main`
is specified as input to fetch `main` branch. The generated changelog will be pushed to the branch.

If you enable [protected-branch][], this action cannot push a commit directly to the branch. Instead,
use `pull_request` input to create a pull request to update the changelog.

```yaml
- uses: rhysd/changelog-from-release/action@v3
  with:
    file: CHANGELOG.md
    github_token: ${{ secrets.GITHUB_TOKEN }}
    pull_request: true
```

Real-world usage example of this action is ['Post release' job of this repository](../.github/workflows/post-release.yml).
Please see [the workflow logs][ci-logs] to know how it runs.

[gh-actions]: https://github.com/features/actions
[release-badge]: https://img.shields.io/github/v/release/rhysd/changelog-from-release.svg
[ci-logs]: https://github.com/rhysd/changelog-from-release/actions/workflows/post-release.yml
[protected-branch]: https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/about-protected-branches
