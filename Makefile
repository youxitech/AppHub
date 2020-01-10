SHELL := /bin/bash
GO_FLAGS := -ldflags="-X 'main.appVersion=$$(git tag --points-at HEAD)' -X 'main.appHash=$$(git rev-parse --short HEAD)'"

run: install
	apphub -d tmp/debug.sqlite3 -r tmp/data --host test.com
.PHONY: run

rundev:
	apphub -p 3389 -d tmp/debug.sqlite3 -r tmp/data --host test.com
.PHONY: rundev

install:
	go install -v $(GO_FLAGS)
.PHONY: install

bindata:
	go-bindata -debug -prefix static static/...
.PHONY: bindata

bindata-prod:
	go-bindata -prefix static static/...
.PHONY: bindata-prod

up:
	sql-migrate up
.PHONY: up

down:
	sql-migrate down
.PHONY: down

# prod build
bundle:
# 	rm -rf static/ui
# 	cd ui && yarn && yarn build
# 	make bindata-prod
	tag=$$(git tag --points-at HEAD) && version=$${tag:-debug} && \
	go build $(GO_FLAGS) -o tmp/apphub-$$version-darwin-amd64 && \
	GOOS=linux GOARCH=amd64 go build $(GO_FLAGS) -o tmp/apphub-$$version-linux-amd64
# 	make bindata
.PHONY: bundle
