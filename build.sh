#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

mkdir -p out

cp scripts/* out/

GOOS=linux   GOARCH=arm64 go build -o out/save_stream.linux-arm64.out   save_stream.go
GOOS=linux   GOARCH=amd64 go build -o out/save_stream.linux-amd64.out   save_stream.go
GOOS=windows GOARCH=amd64 go build -o out/save_stream.windows-amd64.exe save_stream.go
