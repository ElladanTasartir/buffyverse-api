.PHONY: build
build:
	@go build -o ./buffyverse-api ./cmd/buffyverse-api/main.go

.PHONY: run
run:
	@air -c ./build/.air.toml