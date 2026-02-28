# Introduction

`req` is a requirements gathering tool.   It is intended to gather and refine requirements into a specification which can be given to a software factory.

## Beads Setup

```bash
bd init --prefix req
bd doctor
bd sync
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



