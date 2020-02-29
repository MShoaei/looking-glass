#! /bin/sh
for dir in ./*/; do
  cd "$dir" && rm ./*.so; go build -buildmode=plugin . && cd ..
done
