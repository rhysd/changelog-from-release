#!/bin/sh

set -e

# GitHub workspace directory is owned by different user from root. Accessing it is not allowed by default.
# https://github.blog/2022-04-12-git-security-vulnerability-announced/
git config --global --add safe.directory /github/workspace

# Ignore files on LFS (#24)
export GIT_LFS_SKIP_SMUDGE=1

cd "$GITHUB_WORKSPACE" || exit 1

if [ -z "$INPUT_VERSION" ]; then
    TAG_FROM_PAYLOAD="$(jq -r .release.tag_name < "$GITHUB_EVENT_PATH")"
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
echo "::debug::Create pull request?: ${INPUT_PULL_REQUEST}"
echo "::debug::Commit summary template: '${INPUT_COMMIT_SUMMARY_TEMPLATE}'"
echo "::debug::Command arguments: '${INPUT_ARGS}'"
echo "::debug::Header: '${INPUT_HEADER}'"
echo "::debug::Footer: '${INPUT_FOOTER}'"

echo "changelog-from-release version: $(/changelog-from-release -v)"

CHANGELOG="$(GITHUB_TOKEN="$INPUT_GITHUB_TOKEN" /changelog-from-release ${INPUT_ARGS})"

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
    if [ "$INPUT_PULL_REQUEST" = 'true' ]; then
        # Create a branch for pull request
        git checkout -b "changelog-${INPUT_VERSION}"
    fi

    git add "${INPUT_FILE}"
    git \
        -c "user.name=${GITHUB_ACTOR}" \
        -c "user.email=${GITHUB_ACTOR}@users.noreply.github.com" \
        commit -m "${COMMIT_SUMMARY}

    This commit was created by changelog-from-release in '${GITHUB_WORKFLOW}' CI workflow"

    git show

    if [ "$INPUT_PUSH" = 'true' ] || [ "$INPUT_PULL_REQUEST" = 'true' ]; then
        git push --force "https://${GITHUB_ACTOR}:${INPUT_GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git"
    fi
    if [ "$INPUT_PULL_REQUEST" = 'true' ]; then
        echo "${INPUT_GITHUB_TOKEN}" | gh auth login --with-token
        gh pr create --head "changelog-${INPUT_VERSION}" --title "Update changelog for ${INPUT_VERSION}" --body "This PR was automatically created by [changelog-from-release](https://github.com/rhysd/changelog-from-release) action for ${INPUT_VERSION}"
        # Back to the original branch
        git checkout -
    fi
    set +x
fi

echo 'Done'
