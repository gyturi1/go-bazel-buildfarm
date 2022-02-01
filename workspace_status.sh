#!/bin/bash

set -eo pipefail # exit immediately if any command fails.

#repo_url=$(git config --get remote.origin.url")
repo_url=NONE
echo "REPO_URL $repo_url"

#commit_sha=$(git rev-parse HEAD")
commit_sha=NOTCOMMITED
echo "COMMIT_SHA $commit_sha"

#git_branch=$(git rev-parse --abbrev-ref HEAD)
git_branch=MAIN
echo "GIT_BRANCH $git_branch"

#git_tree_status=$(git diff-index --quiet HEAD -- && echo 'Clean' || echo 'Modified')
git_tree_status=NA
echo "GIT_TREE_STATUS $git_tree_status"
