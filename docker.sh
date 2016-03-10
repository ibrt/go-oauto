#!/usr/bin/env bash

PKG_NAME="$(go list)"
BINARY_NAME="${PKG_NAME##*/}"
ROOT_DIR="$(pwd)/../.."

# Create temporary directory.
echo "-> Creating temporary directory..."
TMP_DIR="$(mktemp -d)"
if [[ $? != 0 ]]; then
    echo "Error: Cannot create temporary directory."
    exit 1
fi
trap "rm -rf $TMP_DIR" EXIT
cd "$TMP_DIR"

# Build Go binary.
echo "-> Building Go binary..."
CGO_ENABLED=0 GOOS="linux" GOARCH="amd64" go build -v -a -installsuffix cgo -ldflags "-s" "$PKG_NAME"
if [[ $? != 0 ]]; then
    echo "Error: Cannot build Go binary."
    exit 1
fi

# Build Docker image.
echo "-> Building Docker image..."
cp "$ROOT_DIR/Dockerfile" .
docker build -t "connect-$BINARY_NAME:untested" .
if [[ $? != 0 ]]; then
    err "Error: Cannot build Docker image."
    exit 1
fi