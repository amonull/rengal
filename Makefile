.PHONY: format, report

MAKEFLAGS += --silent

ldflags := -X 'github.com/amonull/rengal/constant.BuiltAt=$(shell date -u)'
ldflags += -X 'github.com/amonull/rengal/constant.BuiltBy=$(shell whoami)'
ldflags += -X 'github.com/amonull/rengal/constant.Revision=$(shell git rev-parse --short HEAD)'
ldflags += -s
ldflags += -w

build_flags := -ldflags=${ldflags}

OUTDIR := ./out

$(shell mkdir -p $(OUTDIR))

all: help

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build        Build the mangal binary"
	@echo "  install      Install the mangal binary"
	@echo "  uninstall    Uninstall the mangal binary"
	@echo "  test         Run the tests"
	@echo "  format 	  formats code using golangci-lint fmt"
	@echo "  report		  runs golangci-lint run and writes to out/"
	@echo "  gif          Generate usage gifs"
	@echo "	 soft-clean   removes out folder"
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

# stores diff first and then runs cmd (could not figure out how to do both at the same time)
format:
	@golangci-lint fmt -d | tee $(OUTDIR)/golangci-lint-fmt.diff
	@golangci-lint fmt

report:
	@golangci-lint run | tee $(OUTDIR)/golangci-lint-report.txt

soft-clean:
	@rm -rf $(OUTDIR)

gif:
	@vhs assets/tui.tape
	@vhs assets/inline.tape
