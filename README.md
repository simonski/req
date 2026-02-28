# Introduction

`req` is a requirements gathering tool.   It is intended to gather and refine requirements into a specification which can be given to a software factory.

## Beads Setup

```bash
# if on an intel mac
go install github.com/steveyegge/beads@v0.47.1

# else
go install github.com/steveyegge/beads@latest

make setup
make tools
make reset

```


See the beads - should be zero beads.   

> Note: if you run `make reset` and you use `VSCode` I advise you to restart VSCode as the beads daemon is a bit flaky.

```bash
bd count
```


## Create requirements

This will generate beads instructions from the requirements document.

```bash
export PATH=./bin:$PATH
parser -f REQUIREMENTS.md > commands.sh
bash commands.sh
```

See the beads - should be 65 beads.

```bash
bd count
```

## see what ralph would do

```bash
wiggum check -name ralph check
```

## Open VSCode and install a beads plugin

Observe the kanban

## simulate a beads loop

```bash
# name is ralph and ralph works fast
# max is 0 (loop till done)
# dryrun means dont really do the work but simulate it
export PATH=./bin:$PATH
wiggum loop -name ralph -max 1 -dryrun 5 -sleep 1

wiggum loop -name ralph -max 1 -dryrun 5 -sleep 1
```



Refresh the kanban

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



