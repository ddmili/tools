#  编译

#!/bin/sh
echo "build...."

#SET CGO_ENABLED=1
#SET GOARCH=
#SET GOOS=windows
#go build main.go

#SET CGO_ENABLED=0
#SET GOOS=darwin
#SET GOARCH=amd64
#go build main.go

#SET CGO_ENABLED=0
#SET GOOS=linux
#SET GOARCH=amd64
#go build main.go

CGO_ENABLED=0 GOOS=windows GOARCH= go build main.go

echo "build success...."