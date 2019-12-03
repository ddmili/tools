#!/bin/sh
echo "build...."

CGO_ENABLED=0 GOOS=windows GOARCH= go build main.go

echo "build success...."