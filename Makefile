BOOTSTRAP_VERSION=3.3.5
JQUERY_VERSION=2.1.4

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

all: $(BOOTSTRAP_FILES) $(JQUERY_FILE) goof

goof: $(wildcard *.go)
	go build

.PHONY: clean
clean:
	git clean -dXf

$(JQUERY_FILE):
	mkdir -p vendor/jquery
	wget --no-verbose "$(JQUERY_URL)" -O "$(JQUERY_FILE)"
	cd vendor && md5sum --check jquery.md5

$(BOOTSTRAP_FILES): extracted_dir=$(basename $(BOOTSTRAP_ARCHIVE))
$(BOOTSTRAP_FILES):
	rm -r "vendor/bootstrap" 2> /dev/null || true
	echo $(extracted_dir)
	wget --no-verbose --no-clobber "$(BOOTSTRAP_URL)"
	mkdir -p vendor
	unzip -d vendor "$(BOOTSTRAP_ARCHIVE)"
	mv "vendor/$(extracted_dir)" "vendor/bootstrap"
	cd vendor && md5sum --check bootstrap.md5
	rm "$(BOOTSTRAP_ARCHIVE)"
