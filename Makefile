.PHONY: dev

dev:
	sass burntPocket/static/static/css
	go generate ./...
	mkdir -p run
	cd run && go run ../burntPocket