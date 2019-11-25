#!/bin/bash
XPKG="github.com/go-x-pkg/sample-app/appversion"
VERSION="0.0.1"
GO_BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%S)
GOOS="linux"
GOARCH="amd64"
GO_LDFLAGS="-s -w -X \"${XPKG}.BuildDate=${GO_BUILD_DATE}\" -X \"${XPKG}.Version=${VERSION}\" -extldflags -static"

mkdir -p ./bin

GOOS=${GOOS} GOARCH=${GOARCH} \
go build \
  -v \
  -ldflags "${GO_LDFLAGS}" -o ./bin/sample-app-server ./cmd/server/main.go
