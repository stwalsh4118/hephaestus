# 1-3 Air Package Guide

Date: 2026-02-18

## Source

- https://github.com/air-verse/air

## Installation API

`go install` (recommended by upstream docs):

```bash
go install github.com/air-verse/air@latest
```

## Usage API

Run with explicit config file:

```bash
air -c .air.toml
```

Initialize default config file:

```bash
air init
```

## Relevant Patterns for This Task

- Install `air` in the backend dev image so the container can watch Go files.
- Keep source mounted into the container (`./backend:/app`) and run `air -c .air.toml`.
- Configure `.air.toml` with:
  - `build.cmd` to compile `./cmd/server`
  - `build.bin` to run compiled binary
  - `include_ext` and `exclude_dir` for efficient watch scope
