.PHONY: ast
ast:
	go build tools/GenerateAst.go
	./GenerateAst ./ast
	go fmt ./ast/...

.PHONY: token
token: 
	go generate ./token

.PHONY: build
build:
	go build -o build/glox

full: token ast build

.PHONY: test
test: build
	go test ./test/...

run:
	./build/glox

dev: build run
