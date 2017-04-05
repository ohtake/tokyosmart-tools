#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

git checkout -B gh-pages

mkdir -p out
GOOS=linux GOARCH=amd64 go build -o out/save_stream.linux-amd64 save_stream.go
GOOS=windows GOARCH=amd64 go build -o out/save_stream.windows-amd64.exe save_stream.go
GOOS=windows GOARCH=386 go build -o out/save_stream.windows-386.exe save_stream.go

git add out -f

if [ "true" = "${TRAVIS}" ]; then
  git config user.name "travis"
  git config user.email "travis@example.net"
fi

git commit -m build
# Needs --quiet to hide authentication token 
git push "${GIT_REMOTE_URL_WITH_AUTH}" gh-pages -f --quiet
git checkout -
