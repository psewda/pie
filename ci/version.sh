#!/bin/bash

build="${TRAVIS_BUILD_NUMBER:-1}"
pkg="github.com/psewda/pie/app"
version="$(cat version) (build $build)"
golang=$(go version | cut -d " " -f 3)
git_commit=$(git rev-parse HEAD)
built=$(date -u)
os_arch=$1

ldflags="-X '$pkg.Version=$version'"
ldflags="$ldflags -X '$pkg.Golang=$golang'"
ldflags="$ldflags -X '$pkg.GitCommit=$git_commit'"
ldflags="$ldflags -X '$pkg.Built=$built'"
ldflags="$ldflags -X '$pkg.OsArch=$os_arch'"

echo $ldflags
