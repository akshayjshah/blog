.DEFAULT_GOAL := build

.PHONY: bootstrap
bootstrap:
	brew install hugo sassc yuicompressor
	# Add GitHub Pages repo for deployment.
	git remote add -f deploy https://github.com/akshayjshah/akshayjshah.github.io.git
	git subtree add --prefix public deploy master --squash


css: styles
	sassc styles/style.scss | yuicompressor --type css > static/style.min.css

.PHONY: clean
clean:
	rm -rf public

build: css archetypes content data layouts static config.toml
	hugo

.PHONY: serve
serve: build
	hugo server --watch

.PHONY: deploy
deploy: build
	git subtree push --prefix=public deploy master
