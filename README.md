goof
====
## Installation
### Prerequisites
[Go](https://golang.org/) and [node-sass](https://github.com/sass/node-sass) are required to run the Makefile.

Versions used during development:
```shell
$ go version
go version go1.5.1 linux/amd64

$ node-sass --version
node-sass	3.3.3	(Wrapper)	[JavaScript]
libsass  	3.2.5	(Sass Compiler)	[C/C++]
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
├── templates       # go html/template files
└── vendor
```
