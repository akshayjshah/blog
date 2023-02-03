.DEFAULT_GOAL := run
SHELL := bash
.SHELLFLAGS := -euo pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

_FLAGS ?= ""
_POSTS := building-a-blog.md \
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
		 grit.md \
		 recipes/waffles.md \
		 recipes/hokkaido-milk-bread.md \
		 recipes/pizza.md \
		 automating-gmail-with-appsscript.md \
		 recipes/coriander-mint-chutney.md \
		 recipes/methi-dal.md \
		 sourdough.md \
		 recipes/pancakes.md \
		 go-time-protobuf-grpc.md \
		 grpc-doesnt-need-trailers.md

.PHONY: help
help: ## Describe useful make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-30s %s\n", $$1, $$2}'

.PHONY: clean
clean: ## Delete intermediate build artifacts
	rm -rf tmp cmd/serve/static cmd/serve/recipes
	rm -f cmd/serve/favicon.ico cmd/serve/*.html

.PHONY: upgrade
upgrade: ## Upgrade Go dependencies
	go get -u -t ./...
	go mod tidy -v

.PHONY: run
run: build tmp/serve ## Run on :8080
	@./tmp/serve

.PHONY: deploy
deploy: INITIAL = $(strip $(shell gcloud config get-value account))
deploy: WANT := akshay@akshayshah.org
deploy: PROJECT := blog-276404/blog
deploy: ## Deploy to GCP
	# Gross, but ensures that I don't forget to run `make clean`.
	$(MAKE) clean
	$(MAKE) build
ifneq ($(WANT),$(INITIAL))
	gcloud config set account $(WANT)
endif
	gcloud auth configure-docker us-central1-docker.pkg.dev
	docker build -t us-central1-docker.pkg.dev/$(PROJECT)/blog:latest .
	docker push us-central1-docker.pkg.dev/$(PROJECT)/blog:latest
	gcloud run deploy blog \
		--image=us-central1-docker.pkg.dev/$(PROJECT)/blog:latest \
		--port=8080 \
		--concurrency=512 \
		--cpu=2 \
		--ingress=all \
		--max-instances=5 \
		--min-instances=default \
		--memory=512Mi \
		--platform=managed \
		--timeout=30s \
		--use-http2 \
		--allow-unauthenticated \
		--cpu-throttling \
		--region=us-central1
ifneq ($(WANT),$(INITIAL))
	gcloud config set account $(INITIAL)
endif

define render-post
	@echo "Rendering $<..."
	@mkdir -p $(@D)
	@tmp/build -style style.css $(_FLAGS) $< > $@
endef

define post-template
build: cmd/serve/$(basename $1).html
cmd/serve/$(basename $1).html: $1 tmp/build page.html style.css
	$$(render-post)
endef

.PHONY: build
build:
	@rm -rf cmd/serve/static && cp -R static cmd/serve
	@cp -p favicon.ico cmd/serve/favicon.ico

tmp/build: go.mod cmd/build/main.go
	@go build -o tmp/build ./cmd/build

tmp/serve: build go.mod cmd/serve/main.go
	@go build -o tmp/serve ./cmd/serve

build: cmd/serve/404.html
cmd/serve/404.html: _FLAGS = -nodates
cmd/serve/404.html: 404.md tmp/build style.css
	$(render-post)

tmp/index.md: index.md style.css tmp/build $(_POSTS)
	@mkdir -p $(@D)
	@tmp/build -recipes recipes -style style.css -index $(_POSTS) > $@

build: cmd/serve/index.html
cmd/serve/index.html: _FLAGS = -nodates -nohome
cmd/serve/index.html: tmp/index.md style.css tmp/build
	$(render-post)

build: cmd/serve/colophon.html
cmd/serve/colophon.html: colophon.md tmp/build page.html style.css
	$(render-post)

build: cmd/serve/license.html
cmd/serve/license.html: _FLAGS = -nolicense -nodates
cmd/serve/license.html: license.md tmp/build page.html style.css
	$(render-post)

$(foreach post,$(_POSTS),$(eval $(call post-template,$(post))))
