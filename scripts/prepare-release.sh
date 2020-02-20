#!/bin/bash

set -e

# Arguments check
if [[ "$#" != 2 ]] || [[ "$1" == '--help' ]]; then
    echo 'Usage: prepare-release.sh {old-version} {new-version}' >&2
    echo '' >&2
    echo "  Release version must be in format 'v{major}.{minor}.{patch}'." >&2
    echo '  After making changes, add --done option and run this script again. It will' >&2
    echo '  push generated tags to remote for release.' >&2
    echo '' >&2
    exit 1
fi

version="$1"
if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo 'Version string in {new-version} argument must match to ''v{major}.{minor}.{patch}'' like v1.2.3' >&2
    exit 1
fi

minor_version="${version%.*}"
major_version="${minor_version%.*}"

prev_version="$2"
if [[ ! "$prev_version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo 'Version string in {old-version} argument must match to ''v{major}.{minor}.{patch}'' like v1.2.3' >&2
    exit 1
fi

# Pre-flight check
if [ ! -d .git ]; then
    echo 'This script must be run at root directory of this repository' >&2
    exit 1
fi

current_branch="$(git symbolic-ref --short HEAD)"
if [[ "$current_branch" != "master" ]]; then
    echo "Release must be created at master branch but current branch is ${$current_branch}" >&2
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
go -v test

# XXX: '.' matches to any character in $prev_version
sed -i '' -E "s/${prev_version}/${pversion}/g" "${files_include_version[@]}"

git add "${files_include_version[@]}"
git commit -m "Bump up version: ${prev_version} -> ${version} [skip action check]"

git tag -d "$major_version" || true
git tag "$major_version"
git tag -d "$minor_version" || true
git tag "$minor_version"
git tag "$version"

# XXX: This script modifies Dockerfile and it tries to fetch a file which does not exist at this point.
# It would cause CI failure.
git push origin "${version}"
git push origin "${minor_version}" --force
git push origin "${major_version}" --force
set +x

if command -v open >/dev/null; then
    open "https://github.com/rhysd/changelog-from-release/releases/new?tag=${version}"
fi

echo "Done."
