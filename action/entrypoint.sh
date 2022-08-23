#!/bin/sh

set -e

# GitHub workspace directory is owned by different user from root. Accessing it is not allowed by default.
# https://github.blog/2022-04-12-git-security-vulnerability-announced/
git config --global --add safe.directory /github/workspace

cd "$GITHUB_WORKSPACE" || exit 1

if [ -z "$INPUT_VERSION" ]; then
    TAG_FROM_PAYLOAD="$(jq .release.tag_name < "$GITHUB_EVENT_PATH")"
    if [ "$TAG_FROM_PAYLOAD" = "null" ]; then
        git fetch --prune --depth=1 origin '+refs/tags/*:refs/tags/*'
        INPUT_VERSION="$(git tag --list | tail -n 1)"
        echo "::debug:: INPUT_VERSION was retrieved from git tags"
    else
        INPUT_VERSION="$TAG_FROM_PAYLOAD"
        echo "::debug:: INPUT_VERSION was retrieved from event payload"
    fi
fi

echo "::debug::Retrieved version: ${INPUT_VERSION}"
echo "::debug::Changelog file: ${INPUT_FILE}"
echo "::debug::Make a commit?: ${INPUT_COMMIT}"
echo "::debug::Push to remote?: ${INPUT_PUSH}"
echo "::debug::Commit summary template: '${INPUT_COMMIT_SUMMARY_TEMPLATE}'"
echo "::debug::Command arguments: '${INPUT_ARGS}'"
echo "::debug::Header: '${INPUT_HEADER}'"
echo "::debug::Footer: '${INPUT_FOOTER}'"

echo "changelog-from-release version: $(/changelog-from-release -v)"

set -x
CHANGELOG="$(GITHUB_TOKEN="$INPUT_GITHUB_TOKEN" /changelog-from-release ${INPUT_ARGS})"
set +x

if [ -n "$INPUT_HEADER" ]; then
    echo "$INPUT_HEADER" > "${INPUT_FILE}"
    echo "$CHANGELOG" >> "${INPUT_FILE}"
else
    echo "$CHANGELOG" > "${INPUT_FILE}"
fi
if [ -n "$INPUT_FOOTER" ]; then
    echo "$INPUT_FOOTER" >> "${INPUT_FILE}"
fi


if [ "$INPUT_COMMIT" = 'true' ]; then
    COMMIT_SUMMARY="Update changelog for ${INPUT_VERSION}"
    if [ -n "$INPUT_COMMIT_SUMMARY_TEMPLATE" ]; then
        COMMIT_SUMMARY="$(printf "$INPUT_COMMIT_SUMMARY_TEMPLATE" "$INPUT_VERSION")"
    fi

    set -x
    git add "${INPUT_FILE}"
    git \
        -c "user.name=${GITHUB_ACTOR}" \
        -c "user.email=${GITHUB_ACTOR}@users.noreply.github.com" \
        commit -m "${COMMIT_SUMMARY}

    This commit was created by changelog-from-release in '${GITHUB_WORKFLOW}' CI workflow"

    if [ "$INPUT_PUSH" = 'true' ]; then
        git push --force "https://${GITHUB_ACTOR}:${INPUT_GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git"
    fi
    set +x
fi

echo 'Done'
