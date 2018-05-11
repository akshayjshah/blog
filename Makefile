.DEFAULT_GOAL := build
SHELL := bash
.SHELLFLAGS := -euo pipefail -c
.ONESHELL:
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

_FLAGS ?= ""
_POSTS = building-a-blog.md \
		 zero-to-code-monkey.md \
		 audible-literacy-filter.md \
		 language-use-on-github.md \
		 soft-deletion-in-django.md \
		 license-your-code.md \
		 testing-django-fields.md \
		 podcasts-for-developers.md \
		 decade-of-cap.md \
		 getting-started-with-go.md

define render-post
	@echo "Rendering $<..."
	@mkdir -p $(@D)
	@bin/build $(_FLAGS) $< > $@
endef

define post-template
build: site/$(basename $1)/index.html
site/$(basename $1)/index.html: $1 bin/build page.html
	$$(render-post)
endef

.PHONY: clean
clean:
	rm -rf site bin tmp

bin/build: go.mod main.go
	@echo "Compiling generator..."
	@go build -o bin/build .

build:
	@echo "Copying static assets..."
	@rm -rf site/{img,docs}
	@cp -R static/* site

build: site/404.html
site/404.html: _FLAGS = -nodates
site/404.html: 404.md bin/build
	$(render-post)

tmp/index.md: index.md bin/build $(_POSTS)
	@echo "Generating index..."
	@mkdir -p $(@D)
	@bin/build -index $(_POSTS) > $@

build: site/index.html
site/index.html: _FLAGS = -nodates -nohome
site/index.html: tmp/index.md bin/build
	$(render-post)

build: site/colophon/index.html
site/colophon/index.html: colophon.md bin/build page.html
	$(render-post)

build: site/license/index.html
site/license/index.html: _FLAGS = -nolicense -nodates
site/license/index.html: license.md bin/build page.html
	$(render-post)

$(foreach post,$(_POSTS),$(eval $(call post-template,$(post))))
