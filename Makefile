.PHONY: dev

dev:
	sass abigit/static/static/css
	go generate ./...
	mkdir -p run
	cd run && go run ../abigit