SHELL=/bin/bash

all: gen_comp

gen_comp:
	go run $(GOPATH)/src/github.com/antha-lang/antha/cmd/antha/antha.go -outdir=lib an
	gofmt -w -s lib

test:
	go test -v `go list ./... | grep -v vendor | grep -v bvendor`

fmt_json:
	for i in `find examples -name '*.json' -o -name '*.yml'`; do \
	  python -mjson.tool "$$i" > "$$i.bak" && mv "$$i.bak" "$$i"; \
	done

compile:
	go install -v github.com/antha-lang/elements/cmd/antharun

test_workflows: compile
	for d in `find examples -type d -o -name '*.yml'`; do \
	  if [[ -f "$$d/workflow.json" && -f "$$d/parameters.yml" ]]; then \
	    /bin/echo -n "Checking $$d..."; \
	    (cd "$$d" && antharun --workflow workflow.json --parameters parameters.yml $(ANTHA_ARGS) > /dev/null); \
	    if [[ $$? == 0 ]]; then \
	      /bin/echo "OK"; \
	    else \
	      /bin/echo "FAIL"; \
	    fi; \
	  fi; \
	done

.PHONY: all gen_comp fmt_json test test_workflows ALWAYS
