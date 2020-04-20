#!/bin/bash

build="${TRAVIS_BUILD_NUMBER:-0}"
version=$(cat version)-$build
linux_archive=pie-$version-linux-amd64.tar.gz
windows_archive=pie-$version-windows-amd64.zip
cd ./bin
mkdir -p package/$version

cd ./linux-amd64 && tar -zcvf $linux_archive pie && cd -
mv -fv ./linux-amd64/$linux_archive ./package/$version

cd ./windows-amd64 && zip $windows_archive pie.exe && cd -
mv -fv ./windows-amd64/$windows_archive ./package/$version
