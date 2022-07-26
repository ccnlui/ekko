# PHONY targets are not really the name of a file; rather it is just a name for
# a recipe to be executed when you make an explicit request

GOPATH := $(shell go env GOPATH)

.PHONY: build
build:
	go build ./...

.PHONY: generate
generate:
	protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./proto/ekko.proto
