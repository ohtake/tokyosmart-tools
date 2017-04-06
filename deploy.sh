#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

git checkout -B gh-pages

./build.sh

git add out -f

if [ "true" = "${TRAVIS}" ]; then
  git config user.name "travis"
  git config user.email "travis@example.net"
fi

git commit -m build
# Needs --quiet to hide authentication token 
git push "${GIT_REMOTE_URL_WITH_AUTH}" gh-pages -f --quiet
git checkout -
