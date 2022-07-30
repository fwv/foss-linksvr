.DEFAULT_GOAL := all

all:
	go build -o bin/linksvr/linksvr cmd/main.go

PROTO_PATHS = $(shell find pkg/ internal/ -name '*.proto' | xargs -I {} dirname {} | uniq)

proto: proto_path

proto_path: $(PROTO_PATHS)
	@$(foreach p,$^,deps/protoc/bin/protoc -I . --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $(wildcard $(p)/*.proto);)