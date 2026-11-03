#!/bin/sh
# Builds the OnionOS package: dist/RabbitEars-v<version>-OnionOS.zip
set -e
cd "$(dirname "$0")"
VERSION=$(grep -o 'const version = "[^"]*"' main.go | cut -d'"' -f2)
[ -f go.work ] || go mod tidy
go test ./...
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -trimpath -ldflags="-s -w" -o package/RabbitEars/rabbitears .
rm -rf dist && mkdir -p dist/stage/App
cp -r package/RabbitEars dist/stage/App/
rm -rf dist/stage/App/RabbitEars/data
cp -r package/Media dist/stage/
(cd dist/stage && zip -qr "../RabbitEars-v$VERSION-OnionOS.zip" App Media)
rm -rf dist/stage
cp converter/vhsify dist/vhsify
(cd converter && zip -q ../dist/vhsify-pi-docker.zip vhsify watch.sh Dockerfile docker-compose.example.yml)
echo "Built dist/RabbitEars-v$VERSION-OnionOS.zip dist/vhsify and dist/vhsify-pi-docker.zip"
