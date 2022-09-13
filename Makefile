.PHONY: dev templates styles

templates:
	go generate ./abigit/http/views

styles:
	sass abigit/static/static/css

dev: templates styles
	mkdir -p run
	cd run && go run ../abigit