# oapit

OpenAPI 3.0 CLI toolkit

## Usage

```prompt
$ oapit -f schema.yml validate MySchema ./payload.json
```

## Install

Pre-built binaries are available on: https://github.com/karupanerura/oapit/releases/tag/v0.0.1

```prompt
$ VERSION=0.0.1
$ curl -sfLO https://github.com/karupanerura/oapit/releases/download/v${VERSION}/oapit_${VERSION}_$(go env GOOS)_$(go env GOARCH).tar.gz
$ tar zxf oapit_${VERSION}_$(go env GOOS)_$(go env GOARCH).tar.gz
$ install -m 0755 oapit $PREFIX
$ rm oapit oapit_${VERSION}_$(go env GOOS)_$(go env GOARCH).tar.gz
```
