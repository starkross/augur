---
id: install
title: Install
sidebar_position: 2
---

# Install

Pick whichever method fits your workflow. Every release ships signed binaries, Docker images, and SBOMs — see [Security](./security) for verification.

## Homebrew

```sh
brew install --cask starkross/tap/augur
```

## Go

```sh
go install github.com/starkross/augur/cmd/augur@latest
```

## Docker

```sh
docker run --rm -v "$(pwd):/work" ghcr.io/starkross/augur:latest config.yaml
```

## Binary releases

Download prebuilt binaries from [GitHub Releases](https://github.com/starkross/augur/releases). Artifacts are available for Linux, macOS, and Windows on both `amd64` and `arm64`.

## Verify installation

```sh
augur --version
```
