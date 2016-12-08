# Antha Elements

[![GoDoc](http://godoc.org/github.com/antha-lang/elements?status.svg)](http://godoc.org/github.com/antha-lang/elements)
[![Build Status](https://travis-ci.org/antha-lang/elements.svg?branch=master)](https://travis-ci.org/antha-lang/elements)

Main instructions are in [antha-lang/antha](https://github.com/antha-lang/antha).

To build or update elements:
```sh
make
# or
make current
```

By default, `make` will download and update any dependent libraries. If you
have any modifications to these dependencies (e.g., non-master branches), `make
current` will build elements without updating any dependent libraries.

To run tests:
```sh
make test
```

To run examples manually:
```sh
make
cd examples/X/Y && antharun
```
