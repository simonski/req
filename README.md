# Introduction

`req` is a requirements gathering tool.   It is intended to gather and refine requirements into a specification which can be given to a software factory.

## Beads Setup

```bash
go install github.com/steveyegge/beads@v0.47.1
bd init --prefix req
bd sync
bd migrate sync beads-sync
bd ready
make tools
```

## Create requirements

This will generate beads instructions from the requirements document.

```bash
export PATH=./bin:$PATH
parser -f REQUIREMENTS.md > commands.sh
bash commands.sh
```

See the beads

```bash
bd list
```

## see what ralph would do

```bash
wiggum check -name ralph check
```

## Open VSCode and install a beads plugin

Observe the kanban

## simulate a beads loop

Open 3 terminals and `export PATH=./bin:$PATH` in each

```bash
# name is ralph and ralph works fast
# max is 0 (loop till done)
# dryrun means dont really do the work but simulate it
wiggum loop -name ralph -max 0 -dryrun -sleep 1
```

```bash
# name is jane and jane works half speed
# max is 0 (loop till done)
# dryrun means dont really do the work but simulate it
wiggum loop -name jane -max 0 -dryrun -sleep 2
```

```bash
# name is manpreet and manpreet takes their time
# max is 0 (loop till done)
# dryrun means dont really do the work but simulate it
wiggum loop -name manpreet -max 0 -dryrun -sleep 4
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



