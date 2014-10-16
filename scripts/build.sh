#!/usr/bin/env bash
set -e

OWNER=ninjablocks
BIN_NAME=sphere-tools
PROJECT_NAME=sphere-tools


# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

GIT_COMMIT="$(git rev-parse HEAD)"
GIT_DIRTY="$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)"
VERSION="$(grep "const Version " version.go | sed -E 's/.*"(.+)"$/\1/' )"

# remove working build
# rm -rf .gopath
if [ ! -d ".gopath" ]; then
	mkdir -p .gopath/src/github.com/${OWNER}
	ln -sf ../../../.. .gopath/src/github.com/${OWNER}/${PROJECT_NAME}
fi

export GOPATH="$(pwd)/.gopath"

if [ ! -d $GOPATH/src/github.com/ninjasphere/go-ninja ]; then
	# Clone our internal commons package
	git clone git@github.com:ninjasphere/go-ninja.git $GOPATH/src/github.com/ninjasphere/go-ninja
fi

# move the working path and build
cd .gopath/src/github.com/${OWNER}/${PROJECT_NAME}
go get -d -v ./...

# build each of the tools and put them in the bin folder for packaging.
go build -ldflags "-X main.GitCommit ${GIT_COMMIT}${GIT_DIRTY}" -o ./bin/sphere-serial -x github.com/ninjasphere/go-ninja/tools/sphere-serial