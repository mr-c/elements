SHELL=/bin/bash

AN_DIRS=an
AN_OUT=lib
INPUT_DIRS=an
PACKAGE=github.com/antha-lang/elements

# Compile after downloading dependencies
all: update_deps fmt_json compile

# Compile using current state of working directories
current: fmt_json compile

gen_comp: anthac anthafmt
	anthafmt -w $(AN_DIRS)
	antha -outdir=$(AN_OUT) $(AN_DIRS)
	gofmt -w -s lib

test: check_json gen_comp
	go test `go list ./... | grep -v vendor`

check: check_json check_elements

check_elements: anthac
	antha -outdir= $(AN_DIRS)

check_json:
	go run cmd/format-json/main.go $(INPUT_DIRS) > /dev/null

fmt_json:
	go run cmd/format-json/main.go -inPlace $(INPUT_DIRS)

update_deps:
	go list -f '{{join .Deps "\n"}}' $(PACKAGE)/cmd/antharun \
	  | grep -v vendor \
	  | grep -v $(PACKAGE) \
	  | xargs go list -f '{{if .Standard}}{{else}}{{.ImportPath}}{{end}}' \
	  | xargs go get -f -u -d -v

anthac:
	go install -v github.com/antha-lang/antha/cmd/antha

anthafmt:
	go install -v github.com/antha-lang/antha/cmd/anthafmt

compile: gen_comp
	go install -v $(PACKAGE)/cmd/antharun
	antharun list elements > /dev/null
