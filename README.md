# Antha Elements

[![GoDoc](http://godoc.org/github.com/antha-lang/elements?status.svg)](http://godoc.org/github.com/antha-lang/elements)
[![Build Status](https://travis-ci.org/antha-lang/elements.svg?branch=master)](https://travis-ci.org/antha-lang/elements)

This repo is for storing and running Antha protocols. 

## Installation
Main instructions are in [antha-lang/antha](https://github.com/antha-lang/antha).


## Build
To build or update elements:
```sh
make
# or
make current
```

or run this command from anywhere:
```sh
make -C "$(go list -f '{{.Dir}}' github.com/antha-lang/elements)"
```


By default, `make` will download and update any dependent libraries. If you
have any modifications to these dependencies (e.g., non-master branches), `make
current` will build elements without updating any dependent libraries.


## Test
To run tests:
```sh
make test
```

To run examples manually:
```sh
make
cd examples/X/Y && antharun
```
## Run 
```sh
antharun --parameters myparameters.json --workflow workflowfile.json
```

## Help
```sh
antharun --help
```

## Academy
Go to the [Antha Academy](https://github.com/antha-lang/elements/tree/master/an/AnthaAcademy) page to be guided through how to use antha in more detail.

