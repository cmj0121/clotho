SRC := $(shell find . -name '*.go')
BIN := $(subst cmd/,dist/,$(wildcard cmd/*))

BIN_PATH := /usr/local/bin

.PHONY: all clean test run build install upgrade help $(SUBDIR)

all: $(SUBDIR) 		# default action
	@[ -f .git/hooks/pre-commit ] || pre-commit install --install-hooks
	@git config commit.template .git-commit-template

clean: $(SUBDIR)	# clean-up environment
	@find . -name '*.sw[po]' -delete
	@rm -f $(BIN)

test:				# run test
	gofmt -w -s $(SRC)
	go test -v ./...

run:				# run in the local environment
	go run

build: $(BIN)		# build the binary/library

install: $(BIN)		# install the binary/library on local
	install -m755 $(BIN) $(BIN_PATH)/

upgrade:			# upgrade all the necessary packages
	pre-commit autoupdate

help:				# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'

dist/%: cmd/%/main.go $(SRC)
	mkdir -p $(@D)
	go build -ldflags="-s -w" -o $@ $<
