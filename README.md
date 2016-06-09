goof
====
## Installation
### Prerequisites
[Go](https://golang.org/) and [npm](https://www.npmjs.com/) are required to run
the Makefile.

Versions used during development:
```shell
$ make version
go    version  go1.5.1   linux/amd64
node  version  v0.10.25
npm   version  1.4.21
```

### Building
```
go get github.com/zenazn/goji
make
```

### Running
Run `./goof` and head to http://localhost:8000/

## Hacking
See the [API](API.md) definition to get an idea of how the JS client works.
