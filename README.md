goof
====
## Installation
### Prerequisites
[Go](https://golang.org/) and [npm](https://www.npmjs.com/) are required to run
the Makefile.

Versions used during development:
```shell
$ go version
go version go1.5.1 linux/amd64
$ node --version
v0.10.25
$ npm --version
1.4.21
```

### Building
Run `make`.

### Running
Run `./goof` and head to http://localhost:8000/

## Files
```
.
├── sass            # unparsed SASS files
├── static          # only directory accessible through HTTP
│   ├── css
│   ├── fonts
│   ├── js
│   └── sass        # compiled SASS
└── templates       # go html/template files
```
