run: install
	apphub
.PHONY: run

install:
	go install -v -ldflags="-X 'main.appVersion=$(git tag --points-at HEAD)' -X 'main.appHash=$(git describe --abbrev=0 --tags)'"
.PHONY: install

bindata:
	go-bindata -debug -prefix static static/...
.PHONY: bindata

up:
	sql-migrate up
.PHONY: bindata
