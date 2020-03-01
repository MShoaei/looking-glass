#! /bin/sh
find . -name "*.so" -exec rm '{}' \;
for dir in ./*/; do
  cd $(git rev-parse --show-toplevel)
  cd "$dir" && go build -buildmode=plugin . && cd ..
done
