quickstart:
	brew install hugo
	npm install -g less less-plugin-clean-css

css:
	lessc --clean-css="--s1 --advanced --compatibility=ie8" less/style.less static/style.min.css

serve:
	hugo server --watch

build:
	hugo
