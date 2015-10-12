BOOTSTRAP_VERSION=3.3.5
JQUERY_VERSION=2.1.4

ifeq ($(shell which node-sass),)
$(error node-sass is not installed)
endif

ifeq ($(shell which go),)
$(error go is not installed)
endif

BOOTSTRAP_ARCHIVE=bootstrap-$(BOOTSTRAP_VERSION)-dist.zip
BOOTSTRAP_URL=https://github.com/twbs/bootstrap/releases/download/v$(BOOTSTRAP_VERSION)/$(BOOTSTRAP_ARCHIVE)
BOOTSTRAP_FILES=vendor/bootstrap/css/bootstrap.min.css vendor/bootstrap/js/bootstrap.min.js\
				vendor/bootstrap/fonts/glyphicons-halflings-regular.eot\
				vendor/bootstrap/fonts/glyphicons-halflings-regular.svg\
				vendor/bootstrap/fonts/glyphicons-halflings-regular.ttf\
				vendor/bootstrap/fonts/glyphicons-halflings-regular.woff\
				vendor/bootstrap/fonts/glyphicons-halflings-regular.woff2\

JQUERY_FILE=vendor/jquery/jquery.min.js
JQUERY_URL=https://code.jquery.com/jquery-$(JQUERY_VERSION).min.js

SASS_SRC=$(shell find sass -type f -name "*.sass")
SASS_COMPILED=$(addsuffix .css,$(addprefix static/,$(basename $(SASS_SRC))))

all: $(BOOTSTRAP_FILES) $(JQUERY_FILE) $(SASS_COMPILED) goof

goof: $(shell find . -type f -name "*.go")
	go build

.PHONY: clean
clean:
	git clean -dXf

.PHONY: watch
watch:
	node-sass --watch --output-style compressed -r sass -o static/sass

$(SASS_COMPILED): $(SASS_SRC)
	node-sass --output-style compressed -r sass -o static/sass

$(JQUERY_FILE):
	mkdir -p vendor/jquery
	wget --no-verbose "$(JQUERY_URL)" -O "$(JQUERY_FILE)"
	cd vendor && md5sum --check jquery.md5

$(BOOTSTRAP_FILES): extracted_dir=$(basename $(BOOTSTRAP_ARCHIVE))
$(BOOTSTRAP_FILES):
	rm -r "vendor/bootstrap" 2> /dev/null || true
	wget --no-verbose --no-clobber "$(BOOTSTRAP_URL)"
	mkdir -p vendor
	unzip -d vendor "$(BOOTSTRAP_ARCHIVE)"
	mv "vendor/$(extracted_dir)" "vendor/bootstrap"
	cd vendor && md5sum --check bootstrap.md5
	rm "$(BOOTSTRAP_ARCHIVE)"
