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
	for tool in $$(find tools -mindepth 2 -maxdepth 2 -type f -name '*.go' | sort); do \
		name=$$(basename $$(dirname $$tool)); \
		printf "Building %s -> bin/%s\n" "$$tool" "$$name"; \
		go build -o "bin/$$name" "$$tool"; \
	done

reset:
	@bd reset --force
	@bd init --prefix req
	@bd sync
	@bd migrate sync beads-sync
	@bd ready

clean:
	@rm -rf bin
