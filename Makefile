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

.PHONY: build
build: css
	hugo

.PHONY: serve
serve: css
	hugo server --watch

.PHONY: deploy
deploy:
	# Always force-push, since the canonical revision history is in this
	# repository.
	git push deploy `git subtree split --prefix=public master`:master --force
