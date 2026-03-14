.PHONY: build build-plugin test clean help

build:
	go build -o bin/golinter ./cmd/golinter

build-plugin:
	go build -buildmode=plugin -o bin/golinter.so ./plugin

test:
	go test ./pkg/analyzer/... -v -race

clean:
	rm -rf bin/

help:
	@echo "build        - собрать линтер"
	@echo "build-plugin - собрать плагин для golangci-lint"
	@echo "test         - запустить тесты"
	@echo "clean        - удалить артефакты"