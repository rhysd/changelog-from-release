#!/bin/bash

set -e

# TODO: Retrieve the {old-version} from `git tag --list | tail -n 1`

# Arguments check
if [[ "$#" != 1 ]] || [[ "$1" == '--help' ]]; then
    echo 'Usage: prepare-release.sh {new-version}' >&2
    echo '' >&2
    echo '  This script creates/updates tags and pushes them to remote for the' >&2
    echo '  new release. Release version value passsed to the first argument' >&2
    echo "  must be in the format 'v{major}.{minor}.{patch}'." >&2
    echo '' >&2
    exit 1
fi

version="$1"
if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo 'Version string in {new-version} argument must match to ''v{major}.{minor}.{patch}'' like v1.2.3' >&2
    exit 1
fi

prev_version="$(git tag --list --sort=version:refname | tail -n 1)"
echo "Bumping the version from ${prev_version} to ${version}"

minor_version="${version%.*}"
major_version="${minor_version%.*}"

# Pre-flight check
if [ ! -d .git ]; then
    echo 'This script must be run at root directory of this repository' >&2
    exit 1
fi

current_branch="$(git symbolic-ref --short HEAD)"
if [[ "$current_branch" != "master" ]]; then
    echo "Release must be created at master branch but current branch is ${current_branch}" >&2
    exit 1
fi

if ! git diff --quiet; then
    echo 'Working tree is dirty! Please ensure all changes are committed and working tree is clean' >&2
    exit 1
fi

if ! git diff --cached --quiet; then
    echo 'Git index is dirty! Please ensure all changes are committed and Git index is clean' >&2
    exit 1
fi

echo "Releasing ${version}... (minor=${minor_version}, major=${major_version})"

files_include_version=( "main.go" "action/Dockerfile" )

set -x
go build
go test
go vet
govulncheck
staticcheck

prev_version_pat="${prev_version//\./\\.}"
sed -i '' -E "s/${prev_version_pat}/${version}/g" "${files_include_version[@]}"
# Replace changelog-from-release_{prev}_{target}.zip -> changelog-from-release_{next}_{target}.zip
sed -i '' -E "s/_${prev_version_pat#v}_/_${version#v}_/g" ./action/Dockerfile

git add "${files_include_version[@]}"
git commit -m "bump version from ${prev_version} to ${version}"
git show HEAD

git tag -d "$major_version" || true
git tag "$major_version"
git tag -d "$minor_version" || true
git tag "$minor_version"
git tag "$version"

git push --no-verify origin master
git push --no-verify origin "${version}"
git push --no-verify origin "${minor_version}" --force
git push --no-verify origin "${major_version}" --force
set +x

echo "Done."
