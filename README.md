# Introduction

`req` is a requirements gathering tool.   It is intended to gather and refine requirements into a specification which can be given to a software factory.

## Beads Setup

```bash
go install github.com/steveyegge/beads@v0.47.1
bd init --prefix req
bd doctor
bd sync
bd migrate sync beads-sync
bd ready
go run parser.go -f REQUIREMENTS.md > commands.sh
bash commands.sh
```

## building

```bash
go run wiggum.go -no-claim -no-branch
```

## Building

```bash
go install github.com/simonski/req@latest
```

## Running

```bash
req init
req server
```

## Usage

Either via the website `http://localhost:8000` or via the terminal using `req`



