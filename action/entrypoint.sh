#!/bin/sh

set -e

cd "$GITHUB_WORKSPACE" || exit 1

if [ -z "$INPUT_VERSION" ]; then
    git fetch --depth=1 origin '+refs/tags/*:refs/tags/*'
    INPUT_VERSION="$(git tag --list | tail -n 1)"
fi

echo "changelog-from-release version: $(/changelog-from-release -v)"

echo "Running changelog-from-release to update ${INPUT_FILE} ..."
GITHUB_TOKEN="$INPUT_GITHUB_TOKEN" /changelog-from-release > "${INPUT_FILE}"

echo "Committing changes in ${INPUT_FILE} ..."
git -c "user.name=${GITHUB_ACTOR}" -c "user.email=${GITHUB_ACTOR}@users.noreply.github.com" add "${INPUT_FILE}"
git -c "user.name=${GITHUB_ACTOR}" -c "user.email=${GITHUB_ACTOR}@users.noreply.github.com" commit -m "Update changelog for ${INPUT_VERSION}"

if [ "$INPUT_PUSH" = 'true' ]; then
    echo 'Pushing generated commit to remote'
    MY_REPO_URL="https://${GITHUB_ACTOR}:${INPUT_GITHUB_TOKEN}@github.com/${REPOSITORY}.git"
    git push --force-with-lease "${MY_REPO_URL}"
fi

echo 'Done'
