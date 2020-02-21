#! /bin/sh
for dir in ./*/; do
  cd "$dir" && go build -buildmode=plugin . && cd ..
done
