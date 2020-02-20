#!/bin/sh

set -e

cd "$GITHUB_WORKSPACE" || exit 1

if [ -z "$INPUT_VERSION" ]; then
    git fetch --prune --depth=1 origin '+refs/tags/*:refs/tags/*'
    INPUT_VERSION="$(git tag --list | tail -n 1)"
fi

echo "::debug::Retrieved version: ${INPUT_VERSION}"
echo "::debug::Make a commit?: ${INPUT_COMMIT}"
echo "::debug::Push to remote?: ${INPUT_PUSH}"

echo "changelog-from-release version: $(/changelog-from-release -v)"

set -x
GITHUB_TOKEN="$INPUT_GITHUB_TOKEN" /changelog-from-release > "${INPUT_FILE}"
set +x

if [ "$INPUT_COMMIT" = 'true' ]; then
    set -x
    git add "${INPUT_FILE}"
    git \
        -c "user.name=${GITHUB_ACTOR}" \
        -c "user.email=${GITHUB_ACTOR}@users.noreply.github.com" \
        commit -m "Update changelog for ${INPUT_VERSION}

    This commit was created by changelog-from-release in '${GITHUB_WORKFLOW}' CI workflow"

    if [ "$INPUT_PUSH" = 'true' ]; then
        git push --force "https://${GITHUB_ACTOR}:${INPUT_GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git"
    fi
    set +x
fi

echo 'Done'
