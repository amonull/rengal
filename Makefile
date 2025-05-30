MAKEFLAGS += --silent

ldflags := -X 'github.com/amonull/rengal/constant.BuiltAt=$(shell date -u)'
ldflags += -X 'github.com/amonull/rengal/constant.BuiltBy=$(shell whoami)'
ldflags += -X 'github.com/amonull/rengal/constant.Revision=$(shell git rev-parse --short HEAD)'
ldflags += -s
ldflags += -w

build_flags := -ldflags=${ldflags}

OUTDIR := ./out

$(shell mkdir $(OUTDIR))

all: help

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build        Build the mangal binary"
	@echo "  install      Install the mangal binary"
	@echo "  uninstall    Uninstall the mangal binary"
	@echo "  test         Run the tests"
	@echo "  clean-code   runs goimport on all go files not in ./vendor"
	@echo "  gif          Generate usage gifs"
	@echo "  help         Show this help message"
	@echo ""

install:
	@go install "$(build_flags)"


build:
	@go build "$(build_flags)"

test:
	@go test ./...

uninstall:
	@rm -f $(shell which mangal)

clean-code:
	@find . -not -path "./vendor/*" -type f -name "*.go" -exec goimports -w -d {} \; > $(OUTDIR)/goimport.diff

gif:
	@vhs assets/tui.tape
	@vhs assets/inline.tape
