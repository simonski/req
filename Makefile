.PHONY: help default tools clean

default: help

help:
	@printf "Available targets:\n\n"
	@printf "  make help    Print this usage message.\n"
	@printf "  make tools   Build all Go tools under tools/ into ./bin.\n"
	@printf "  make clean   Remove built binaries from ./bin.\n"
	@printf "  make reset   Reset beads completely.\n"
	@printf "\n"

tools:
	@mkdir -p bin
	@set -e; \
	for tool in $$(find tools -mindepth 2 -maxdepth 2 -type f -name '*.go' ! -name '*_test.go' | sort); do \
		name=$$(basename $$(dirname $$tool)); \
		printf "Building %s -> bin/%s\n" "$$tool" "$$name"; \
		go build -o "bin/$$name" "$$tool"; \
	done

setup:
	@bd init --prefix req
	@bd sync
	@bd migrate sync beads-sync
	@bd ready
	@echo Restart VSCode

reset:
	@bd list -n 0 --json | jq '.[].id' | xargs bd delete $1 -f
	@bd list -s closed -n 0 --json | jq '.[].id' | xargs bd delete $1 -f

clean:
	@rm -rf bin
