#!/bin/sh

VERSION=$(git tag -l | tail -n 1)
COMMIT=$(git rev-parse --short HEAD)

go run -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT}" cmd/main.go
