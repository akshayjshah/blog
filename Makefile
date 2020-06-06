.DEFAULT_GOAL := build
SHELL := bash
.SHELLFLAGS := -euo pipefail -c
.ONESHELL:
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

_INITIAL_ACCOUNT = $(strip $(shell gcloud config get-value account))
_ACCOUNT = akshay@akshayshah.org
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
		 getting-started-with-go.md \
		 lazy-loading-data-with-swiftui-and-combine.md \
		 recipes/garam-masala.md \
		 recipes/bakers-percentages.md \
		 recipes/sourdough-starter.md \
		 recipes/focaccia.md \
		 recipes/chicken-kebabs.md \
		 recipes/methi-murgh.md \
		 books/how-to-cook-everything.md \
		 books/grit.md \
		 recipes/waffles.md \
		 recipes/hokkaido-milk-bread.md \
		 recipes/pizza.md \
		 automating-gmail-with-appscript.md

define render-post
	@echo "Rendering $<..."
	@mkdir -p $(@D)
	@bin/build -style style.css $(_FLAGS) $< > $@
endef

define post-template
build: site/$(basename $1)/index.html
site/$(basename $1)/index.html: $1 bin/build page.html style.css
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
	@cp favicon.ico site/favicon.ico

build: site/404.html
site/404.html: _FLAGS = -nodates
site/404.html: 404.md bin/build style.css
	$(render-post)

tmp/index.md: index.md style.css bin/build $(_POSTS)
	@echo "Generating index..."
	@mkdir -p $(@D)
	@bin/build -books books -recipes recipes -style style.css -index $(_POSTS) > $@

build: site/index.html
site/index.html: _FLAGS = -nodates -nohome
site/index.html: tmp/index.md style.css bin/build
	$(render-post)

build: site/colophon/index.html
site/colophon/index.html: colophon.md bin/build page.html style.css
	$(render-post)

build: site/license/index.html
site/license/index.html: _FLAGS = -nolicense -nodates
site/license/index.html: license.md bin/build page.html style.css
	$(render-post)

$(foreach post,$(_POSTS),$(eval $(call post-template,$(post))))

.PHONY: deploy
deploy: build
ifneq ($(_ACCOUNT),$(_INITIAL_ACCOUNT))
	gcloud config set account $(_ACCOUNT)
endif
	gsutil -o "GSUtil:parallel_process_count=4" -o "GSUtil:parallel_thread_count=1" -m rsync -d -j -R site gs://www.akshayshah.org
ifneq ($(_ACCOUNT),$(_INITIAL_ACCOUNT))
	gcloud config set account $(_INITIAL_ACCOUNT)
endif
