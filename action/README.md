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
      - uses: actions/checkout@v2
        with:
          ref: master
      - uses: rhysd/changelog-from-release/action@v2
        with:
          file: CHANGELOG.md
          github_token: ${{ secrets.GITHUB_TOKEN }}
```

`file` is a file path for changelog and `github_token` is a GitHub API token to retrieve releases.
Please read [action.yml](action.yml) for more details.

Note that `actions/checkout@v2` does not fetch branches by default. In above example, `ref: master`
is specified as input to update `master` branch.

And ['Post release' job of this repository](../.github/workflows/post-release.yml) is a real-world
usage example of the action.

[gh-actions]: https://github.com/features/actions
[release-badge]: https://img.shields.io/github/v/release/rhysd/changelog-from-release.svg
