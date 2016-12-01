SHELL=/bin/bash

all: fmt_json compile

gen_comp:
	go run $(GOPATH)/src/github.com/antha-lang/antha/cmd/antha/antha.go -outdir=lib an starter
	gofmt -w -s lib

test: check_json gen_comp
	go test `go list ./... | grep -v vendor`

check: check_json check_elements

check_elements:
	go run $(GOPATH)/src/github.com/antha-lang/antha/cmd/antha/antha.go -outdir= an starter

check_json:
	go run cmd/format-json/main.go workflows starter > /dev/null

fmt_json:
	go run cmd/format-json/main.go -inPlace workflows starter defaultparameters

compile: gen_comp
	go install -v github.com/antha-lang/elements/cmd/antharun
	antharun list > /dev/null
