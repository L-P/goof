PREREQUISITES = go npm
_ := $(foreach exec,$(PREREQUISITES),\
	$(if $(shell which $(exec)),_,$(error "$(exec) not found in $$PATH.")))

# This is needed to make JS source maps work. Firefox does not seem to
# understand relative URLs so this will have to do for now.
DEV_URL="http://localhost:3000"

# Order matters.
GOOF_JS=\
static/js/goof/Event.js\
static/js/goof/Calendar.js\
static/js/goof/NextEventsView.js\
static/js/goof/main.js\

# Go packages we define.
PACKAGES=calendar gui

BOWER=node_modules/bower/bin/bower
SASS=node_modules/node-sass/bin/node-sass
UGLIFYJS=node_modules/uglify-js/bin/uglifyjs

NODE_MODULES=$(SASS) $(BOWER) $(UGLIFYJS)

BACKBONE_JS=bower_components/backbone/backbone.js
BOOTSTRAP_JS=bower_components/bootstrap/dist/js/bootstrap.js
JQUERY_JS=bower_components/jquery/dist/jquery.js
MUSTACHE_JS=bower_components/mustache.js/mustache.js
UNDERSCORE_JS=bower_components/underscore/underscore.js
FULLCALENDAR_JS=bower_components/fullcalendar/dist/fullcalendar.js

BOOTSTRAP_FILES=$(BOOTSTRAP_JS)\
bower_components/bootstrap/dist/css/bootstrap.min.css\
bower_components/bootstrap/dist/fonts/glyphicons-halflings-regular.eot\
bower_components/bootstrap/dist/fonts/glyphicons-halflings-regular.svg\
bower_components/bootstrap/dist/fonts/glyphicons-halflings-regular.ttf\
bower_components/bootstrap/dist/fonts/glyphicons-halflings-regular.woff\
bower_components/bootstrap/dist/fonts/glyphicons-halflings-regular.woff2\

BOWER_COMPONENTS=$(JQUERY_JS) $(BOOTSTRAP_FILES) $(UNDERSCORE_JS) $(BACKBONE_JS) $(MUSTACHE_JS) $(FULLCALENDAR_JS)

SASS_SRC=$(shell find sass -type f -name "*.sass")
SASS_COMPILED=$(addsuffix .css,$(addprefix static/,$(basename $(SASS_SRC))))

VENDOR_JS_COMPILED=static/js/vendors.min.js
GOOF_JS_COMPILED=static/js/goof.min.js

all: $(NODE_MODULES) $(BOWER_COMPONENTS) $(SASS_COMPILED) $(VENDOR_JS_COMPILED) $(GOOF_JS_COMPILED) goof

goof: goof.go $(shell find $(PACKAGES) -type f -name "*.go")
	go build

# Order matters.
$(VENDOR_JS_COMPILED): $(JQUERY_JS) $(UNDERSCORE_JS) $(BACKBONE_JS) $(BOOTSTRAP_JS) $(MUSTACHE_JS) $(FULLCALENDAR_JS)
	$(UGLIFYJS) --screw-ie8 --mangle --compress --output "$@" -- $^

$(GOOF_JS_COMPILED): $(GOOF_JS)
	$(UGLIFYJS) --screw-ie8 --mangle --compress --output "$@"\
		--source-map "$@.map" --source-map-url "http://localhost:3000/$@.map"\
		--prefix 2\
		-- $^

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
