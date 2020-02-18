#!/bin/sh

set -e

cd "$GITHUB_WORKSPACE" || exit 1

if [ -z "$INPUT_VERSION" ]; then
    INPUT_VERSION="$(git describe --tags)"
fi

echo "changelog-from-release version: $(/changelog-from-release -v)"

echo "Running changelog-from-release to update ${INPUT_FILE} ..."
GITHUB_TOKEN="$INPUT_GITHUB_TOKEN" /changelog-from-release > "${INPUT_FILE}"

echo "Committing changes in ${INPUT_FILE} ..."
git add "${INPUT_FILE}"
git commit -m "Update changelog for ${INPUT_VERSION}"

if [ "$INPUT_PUSH" = 'true' ]; then
    echo 'Pushing generated commit to remote'
    MY_REPO_URL="https://${GITHUB_ACTOR}:${INPUT_GITHUB_TOKEN}@github.com/${REPOSITORY}.git"
    git push --force-with-lease "${MY_REPO_URL}"
fi

echo 'Done'
