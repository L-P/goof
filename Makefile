PREREQUISITES = go npm
_ := $(foreach exec,$(PREREQUISITES),\
	$(if $(shell which $(exec)),_,$(error "$(exec) not found in $$PATH.")))

SASS=node_modules/node-sass/bin/node-sass
BOWER=node_modules/bower/bin/bower

NODE_MODULES=$(SASS) $(BOWER)

JQUERY_FILE=bower_components/jquery/dist/jquery.min.js

BOOTSTRAP_FILES=\
bower_components/bootstrap/dist/css/bootstrap.min.css\
bower_components/bootstrap/dist/js/bootstrap.min.js\
bower_components/bootstrap/dist/fonts/glyphicons-halflings-regular.eot\
bower_components/bootstrap/dist/fonts/glyphicons-halflings-regular.svg\
bower_components/bootstrap/dist/fonts/glyphicons-halflings-regular.ttf\
bower_components/bootstrap/dist/fonts/glyphicons-halflings-regular.woff\
bower_components/bootstrap/dist/fonts/glyphicons-halflings-regular.woff2\

BOWER_COMPONENTS=$(JQUERY_FILE) $(BOOTSTRAP_FILES)

SASS_SRC=$(shell find sass -type f -name "*.sass")
SASS_COMPILED=$(addsuffix .css,$(addprefix static/,$(basename $(SASS_SRC))))

all: $(NODE_MODULES) $(BOWER_COMPONENTS) $(SASS_COMPILED) goof

goof: $(shell find . -type f -name "*.go")
	go build

.PHONY: clean watch version
clean:
	git clean -dXf

watch:
	$(SASS) --watch --output-style compressed -r sass -o static/sass

version:
	@( \
		go version \
		&& echo -n "node version " && node --version \
		&& echo -n "npm version " && npm --version \
	) | column -t

$(SASS_COMPILED): $(SASS_SRC)
	$(SASS) --output-style compressed -r sass -o static/sass

$(NODE_MODULES):
	npm install

$(BOWER_COMPONENTS):
	$(BOWER) install
