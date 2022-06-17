#!/bin/bash
set -e
echo "Tests"
go test -v -race ./...

echo "Build"
go build .

echo "Integration test"
cat privatetestdata/sample1.modfile | ./godephunter --find="github.com/Route42/golang-commons/httpclient@v0.0.0-20210615"
go mod graph | ./godephunter --find='github.com/davecgh/go-spew@v1.1.1'

echo "Copying to local bin folder"
sudo cp godephunter /usr/local/bin/

echo "DONE :)"
