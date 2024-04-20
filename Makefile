# Based on the Example from Joel Homes, author of "Shipping Go" at
# https://github.com/holmes89/hello-api/blob/main/ch10/Makefile

#  https://stackoverflow.com/a/54776239
SHELL := /bin/bash
GO_VERSION := 1.22  # <1>
COVERAGE_AMT := 70  # should be 80
HEREGOPATH := $(shell go env GOPATH)
CURDIR := $(shell pwd)

# setup: # <2>
# 	install-go
# 	init-go
# 
# install-go: # <3>
# 	wget "https://golang.org/dl/go$(GO_VERSION).linux-amd64.tar.gz"
# 	sudo tar -C /usr/local -xzf go$(GO_VERSION).linux-amd64.tar.gz
# 	rm go$(GO_VERSION).linux-amd64.tar.gz
# 
# init-go: # <4>
#     echo 'export PATH=$$PATH:/usr/local/go/bin' >> $${HOME}/.bashrc
#     echo 'export PATH=$$PATH:$${HOME}/go/bin' >> $${HOME}/.bashrc
# 
# upgrade-go: # <5>
# 	sudo rm -rf /usr/bin/go
# 	wget "https://golang.org/dl/go$(GO_VERSION).linux-amd64.tar.gz"
# 	sudo tar -C /usr/local -xzf go$(GO_VERSION).linux-amd64.tar.gz
# 	rm go$(GO_VERSION).linux-amd64.tar.gz

# update-javascript:
# 	# htmx
# 	wget -o web/static/htmx.min.js https://unpkg.com/htmx.org@1.9.6/dist/htmx.min.js
# 	# hyperscript (as at 15 October 2023)
# 	wget -o web/static/htmx.min.js https://unpkg.com/hyperscript.org@0.9.11

build-cli:
	cd ${CURDIR}/cmd/cli && go test . && echo "---ok---" && go build -o ${CURDIR}/bin/cli .

build-web:
	cd ${CURDIR}/cmd/web && go test . && echo "---ok---" && go build -o ${CURDIR}/bin/webserver .

build-console:
	# @echo ${CURDIR}
	cd ${CURDIR}/cmd/console && go test . && echo "---ok---" && go build -o ${CURDIR}/bin/console .

build-all: build-cli build-web build-console

build-many:
	go test ./... || exit 1
	cd ${CURDIR}/cmd/web     && go test . && ${CURDIR}/bin/builder.sh . webserver
	cd ${CURDIR}/cmd/cli     && go test . && ${CURDIR}/bin/builder.sh . cli
	cd ${CURDIR}/cmd/console && go test . && ${CURDIR}/bin/builder.sh . console

del-bin:
	rm `ls bin/* | grep -v builder`

# build-dev:
# 	go test ./... && echo "---ok---" && go build -o timeaway -tags=development cmd/main.go

test:
	go test ./... -coverprofile=coverage.out
	cd ${CURDIR}/cmd/web     && go test .  -coverprofile=coverage.out
	cd ${CURDIR}/cmd/cli     && go test .  -coverprofile=coverage.out
	cd ${CURDIR}/cmd/console && go test .  -coverprofile=coverage.out

coverage-verbose:
	go tool cover -func coverage.out | tee cover.rpt
	cd ${CURDIR}/cmd/web     && go tool cover -func coverage.out | tee cover.rpt
	cd ${CURDIR}/cmd/cli     && go tool cover -func coverage.out | tee cover.rpt
	# skip console coverage

coverage-ok:
	cat cover.rpt | grep "total:" | awk '{print ((int($$3) > ${COVERAGE_AMT}) != 1) }'

cover-report:
	# this is for the main module only
	go tool cover -html=coverage.out -o cover.html

clean:
	rm $$(find . -name "*cover*html" -or -name "*cover.rpt" -or -name "*coverage.out")

check: check-format check-vet test coverage-verbose coverage-ok cover-report lint 

check-format: 
	test -z $$(go fmt ./...)
	test -z $$(go fmt cmd/web/*go)
	test -z $$(go fmt cmd/cli/*go)
	test -z $$(go fmt cmd/console/*go)

check-vet: 
	test -z $$(go vet ./...)
	cd ${CURDIR}/cmd/web     && test -z $$(go vet .)
	cd ${CURDIR}/cmd/cli     && test -z $$(go vet .)
	cd ${CURDIR}/cmd/console && test -z $$(go vet .)

testme:
	echo $(HEREGOPATH)

install-lint:
	# https://golangci-lint.run/usage/install/#local-installation to GOPATH
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(HEREGOPATH)/bin v1.57.2
	# report version
	${HEREGOPATH}/bin/golangci-lint --version

lint:
	# golangci-lint run -v ./... 
	${HEREGOPATH}/bin/golangci-lint run ./... 
	cd ${CURDIR}/cmd/web     && ${HEREGOPATH}/bin/golangci-lint run .
	cd ${CURDIR}/cmd/cli     && ${HEREGOPATH}/bin/golangci-lint run .
	cd ${CURDIR}/cmd/console && ${HEREGOPATH}/bin/golangci-lint run .

module-update-tidy:
	go get -u ./...
	go mod tidy
	cd ${CURDIR}/cmd/web     && go mod tidy
	cd ${CURDIR}/cmd/cli     && go mod tidy
	cd ${CURDIR}/cmd/console && go mod tidy

