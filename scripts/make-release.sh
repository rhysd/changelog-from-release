#! /bin/bash

set -e

name="changelog-from-release"
command="${name}"
if [[ "$bin" == *windows* ]]; then
    command="${command}.exe"
fi

rm -rf release
gox -arch 'amd64' -os 'linux darwin windows freebsd openbsd netbsd' ./
mkdir -p release
mv ${name}_* release/
cd release
for bin in *; do
    mv "$bin" "$command"
    zip "${bin}.zip" "$command"
    rm "$command"
done
