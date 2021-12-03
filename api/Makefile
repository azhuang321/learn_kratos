GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

.PHONY: user
user:
	protoc --proto_path=. \
		--proto_path=../third_party \
		--go_out=paths=source_relative:. \
		--go-http_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:. \
		--go-errors_out=paths=source_relative:. \
	    $(shell find user -name *.proto)

.PHONY: group
group:
	protoc --proto_path=. \
		--proto_path=../third_party \
		--go_out=paths=source_relative:. \
		--go-http_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:. \
		--go-errors_out=paths=source_relative:. \
	    $(shell find group -name *.proto)