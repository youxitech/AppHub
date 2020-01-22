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

build-ui:
	rm -rf static/ui
	cd ui && yarn && yarn build
.PHONY: build-ui

build-darwin:
	make bindata-prod
	tag=$$(git tag --points-at HEAD) && version=$${tag:-debug} && \
	go build $(GO_FLAGS) -o tmp/apphub-$$version-amd64-darwin
	make bindata
.PHONY: build-darwin

build-linux:
	make bindata-prod
	tag=$$(git tag --points-at HEAD) && hash=$$(git rev-parse --short HEAD) && version=$${tag:-$$hash} && \
	CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build $(GO_FLAGS) -ldflags "-linkmode external -extldflags -static" -o tmp/apphub-$$version-amd64-linux
	make bindata
.PHONY: build-linux

bundle: build-ui build-darwin build-linux
.PHONY: bundle

deploy:
	rsync --update --progress tmp/*linux systatic:/data/apphub/bin/
	ssh systatic systemctl restart apphub
