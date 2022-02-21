#!/bin/bash

# win 32
# GOARCH=386 GOOS=windows go build -o ./bin/classroom-core-win32-x32.exe

# win 64
GOARCH=amd64 GOOS=windows go build -o ./bin/classroom-core-win32-x64.exe

# arm 32
# GOARCH=arm GOOS=linux go build -o ./bin/classroom-core-arm-x32

# arm 64 linux
GOARCH=arm64 GOOS=linux go build -o ./bin/classroom-core-arm-linux-x64

# arm 64 mac
GOARCH=arm64 GOOS=darwin go build -o ./bin/classroom-core-arm-mac-x64

# mac 64
GOARCH=amd64 GOOS=darwin go build -o ./bin/classroom-core-darwin-x64

# linux 64
GOARCH=amd64 GOOS=linux go build -o ./bin/classroom-core-linux-x64