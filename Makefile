run: install
	apphub -d tmp/debug.sqlite3 -r tmp/data --host test.com
.PHONY: run

rundev:
	apphub -p 3389 -d tmp/debug.sqlite3 -r tmp/data --host test.com
.PHONY: rundev

install:
	go install -v -ldflags="-X 'main.appVersion=$(git tag --points-at HEAD)' -X 'main.appHash=$(git describe --abbrev=0 --tags)'"
.PHONY: install

bindata:
	go-bindata -debug -prefix static static/...
.PHONY: bindata

up:
	sql-migrate up
.PHONY: up

down:
	sql-migrate down
.PHONY: down
